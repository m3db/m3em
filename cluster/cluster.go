// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cluster

import (
	"fmt"
	"sync"

	"github.com/m3db/m3cluster/services"
	env "github.com/m3db/m3em/environment"
	"github.com/m3db/m3x/errors"
	"github.com/m3db/m3x/log"
	"github.com/m3db/m3x/sync"
)

var (
	errInsufficientCapacity            = fmt.Errorf("insufficient node capacity in environment")
	errNodeAlreadyUsed                 = fmt.Errorf("unable to add node, already in use")
	errNodeNotInUse                    = fmt.Errorf("unable to remove node, not in use")
	errClusterAlreadySetup             = fmt.Errorf("cluster already setup")
	errClusterUnableToInitialize       = fmt.Errorf("unable to initialize cluster, it needs to be setup")
	errClusterUnableToAlterPlacement   = fmt.Errorf("unable to alter cluster placement, it needs to be initialized/running")
	errUnableToStartUnitializedCluster = fmt.Errorf("unable to start unitialized cluster")
	errClusterUnableToTeardown         = fmt.Errorf("unable to teardown cluster, it has not been setup")
	errUnableToStopNotRunningCluster   = fmt.Errorf("unable to stop cluster, it is running")
)

type idToNodeMap map[string]env.ServiceNode

type svcCluster struct {
	logger       xlog.Logger
	opts         Options
	knownNodes   env.ServiceNodes
	usedNodes    idToNodeMap
	spares       []env.ServiceNode
	placementSvc services.PlacementService

	statusLock sync.RWMutex
	status     Status
	lastErr    error
}

// New returns a new cluster backed by provided service nodes
func New(
	nodes env.ServiceNodes,
	opts Options,
) (Cluster, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return &svcCluster{
		logger:       opts.InstrumentOptions().Logger(),
		opts:         opts,
		knownNodes:   nodes,
		usedNodes:    make(idToNodeMap, len(nodes)),
		placementSvc: opts.PlacementService(),
		status:       ClusterStatusUninitialized,
	}, nil
}

// TODO(prateek): reset initial placement after teardown
// TODO(prateek): use concurrency and other options

type concurrentNodeFn func(env.ServiceNode)

type concurrentNodeExecutor struct {
	wg      sync.WaitGroup
	nodes   env.ServiceNodes
	workers xsync.WorkerPool
	fn      concurrentNodeFn
}

func newConcurrentNodeExecutor(
	nodes env.ServiceNodes,
	concurrency int,
	fn concurrentNodeFn,
) *concurrentNodeExecutor {
	workerPool := xsync.NewWorkerPool(concurrency)
	workerPool.Init()
	return &concurrentNodeExecutor{
		nodes:   nodes,
		workers: workerPool,
		fn:      fn,
	}
}

func (e *concurrentNodeExecutor) run() {
	e.wg.Add(len(e.nodes))
	for idx := range e.nodes {
		node := e.nodes[idx]
		e.workers.Go(func() {
			defer e.wg.Done()
			e.fn(node)
		})
	}
	e.wg.Wait()
}

func (c *svcCluster) Setup() error {
	if c.Status() != ClusterStatusUninitialized {
		return errClusterAlreadySetup
	}

	psvc := c.placementSvc
	if _, _, err := psvc.Placement(); err != nil { // attempt to retrieve current placement
		c.logger.Infof("unable to retrieve current placement, skipping delete attempt")
	} else if err = psvc.Delete(); err != nil { // delete existing placement
		// only logging error instead of letting it flow through as the placement should
		c.logger.Infof("placement deletion during cluster Setup failed: %v", err)
	}

	var (
		svcBuild        = c.opts.ServiceBuild()
		svcConf         = c.opts.ServiceConfig()
		sessionToken    = c.opts.SessionToken()
		sessionOverride = c.opts.SessionOverride()
		lock            sync.Mutex
		multiErr        xerrors.MultiError
	)

	// setup the nodes, in parallel
	executor := newConcurrentNodeExecutor(c.knownNodes, c.opts.NodeConcurrency(), func(node env.ServiceNode) {
		if err := node.Setup(svcBuild, svcConf, sessionToken, sessionOverride); err != nil {
			lock.Lock()
			multiErr = multiErr.Add(err)
			lock.Unlock()
			return
		}
		lock.Lock()
		c.spares = append(c.spares, node)
		lock.Unlock()
	})
	executor.run()

	return c.markStatus(ClusterStatusSetup, multiErr.FinalError())
}

func (c *svcCluster) Initialize(numNodes int) ([]env.ServiceNode, error) {
	if c.Status() != ClusterStatusSetup {
		return nil, errClusterUnableToInitialize
	}

	numSpares := len(c.spares)
	if numSpares < numNodes {
		return nil, errInsufficientCapacity
	}

	usedNodes := c.spares[:numNodes]
	nodes := make([]services.PlacementInstance, 0, numNodes)
	for _, node := range usedNodes {
		c.usedNodes[node.ID()] = node
		nodes = append(nodes, node)
	}
	c.spares = c.spares[numNodes:]

	psvc := c.placementSvc
	_, err := psvc.BuildInitialPlacement(nodes, c.opts.NumShards(), c.opts.Replication())
	if err != nil {
		return nil, err
	}

	c.statusLock.Lock()
	c.status = ClusterStatusInitialized
	c.statusLock.Unlock()
	return usedNodes, nil
}

func (c *svcCluster) AddNode() (env.ServiceNode, error) {
	if status := c.Status(); status != ClusterStatusRunning && status != ClusterStatusInitialized {
		return nil, errClusterUnableToAlterPlacement
	}

	numSpares := len(c.spares)
	if numSpares < 1 {
		return nil, errInsufficientCapacity
	}
	node := c.spares[0]
	psvc := c.placementSvc
	_, _, err := psvc.AddInstance([]services.PlacementInstance{node})
	if err != nil {
		return nil, err
	}
	c.usedNodes[node.ID()] = node
	c.spares = c.spares[1:]
	return node, nil
}

func (c *svcCluster) RemoveNode(i env.ServiceNode) error {
	if status := c.Status(); status != ClusterStatusRunning && status != ClusterStatusInitialized {
		return errClusterUnableToAlterPlacement
	}

	if _, ok := c.usedNodes[i.ID()]; !ok {
		return errNodeNotInUse
	}
	psvc := c.placementSvc
	_, err := psvc.RemoveInstance(i.ID())
	if err != nil {
		return err
	}
	delete(c.usedNodes, i.ID())
	c.spares = append(c.spares, i)
	return nil
}

func (c *svcCluster) ReplaceNode(oldNode env.ServiceNode) (env.ServiceNode, error) {
	if status := c.Status(); status != ClusterStatusRunning && status != ClusterStatusInitialized {
		return nil, errClusterUnableToAlterPlacement
	}

	if _, ok := c.usedNodes[oldNode.ID()]; !ok {
		return nil, errNodeNotInUse
	}

	numSpares := len(c.spares)
	if numSpares < 1 {
		return nil, errInsufficientCapacity
	}

	spareNode := c.spares[0]
	psvc := c.placementSvc
	_, _, err := psvc.ReplaceInstance(oldNode.ID(), []services.PlacementInstance{spareNode})
	if err != nil {
		return nil, err
	}

	delete(c.usedNodes, oldNode.ID())
	c.usedNodes[spareNode.ID()] = spareNode
	c.spares = append(c.spares[1:], oldNode)

	return spareNode, nil
}

func (c *svcCluster) Spares() []env.ServiceNode {
	return c.spares
}

func (c *svcCluster) markStatus(status Status, err error) error {
	c.statusLock.Lock()
	defer c.statusLock.Unlock()

	if err == nil {
		c.status = status
		return nil
	}

	c.status = ClusterStatusError
	c.lastErr = err
	return err
}

func (c *svcCluster) Teardown() error {
	if status := c.Status(); status == ClusterStatusUninitialized {
		return errClusterUnableToTeardown
	}

	var (
		lock     sync.Mutex
		multiErr xerrors.MultiError
	)

	executor := newConcurrentNodeExecutor(c.knownNodes, c.opts.NodeConcurrency(), func(node env.ServiceNode) {
		if err := node.Teardown(); err != nil {
			lock.Lock()
			multiErr = multiErr.Add(err)
			lock.Unlock()
			return
		}
		lock.Lock()
		c.spares = append(c.spares, node)
		lock.Unlock()
	})
	executor.run()

	return c.markStatus(ClusterStatusUninitialized, multiErr.FinalError())
}

func (c *svcCluster) StartInitialized() error {
	if status := c.Status(); status != ClusterStatusInitialized {
		return errUnableToStartUnitializedCluster
	}

	var (
		lock     sync.Mutex
		multiErr xerrors.MultiError
		nodes    = make(env.ServiceNodes, 0, len(c.usedNodes))
	)
	for _, node := range c.usedNodes {
		nodes = append(nodes, node)
	}

	executor := newConcurrentNodeExecutor(nodes, c.opts.NodeConcurrency(), func(node env.ServiceNode) {
		if err := node.Start(); err != nil {
			lock.Lock()
			multiErr = multiErr.Add(err)
			lock.Unlock()
			return
		}
		lock.Lock()
		c.spares = append(c.spares, node)
		lock.Unlock()
	})
	executor.run()

	return c.markStatus(ClusterStatusRunning, multiErr.FinalError())
}

func (c *svcCluster) Start() error {
	if status := c.Status(); status != ClusterStatusInitialized {
		return errUnableToStartUnitializedCluster
	}

	var (
		lock     sync.Mutex
		multiErr xerrors.MultiError
	)

	executor := newConcurrentNodeExecutor(c.knownNodes, c.opts.NodeConcurrency(), func(node env.ServiceNode) {
		if err := node.Start(); err != nil {
			lock.Lock()
			multiErr = multiErr.Add(err)
			lock.Unlock()
			return
		}
		lock.Lock()
		c.spares = append(c.spares, node)
		lock.Unlock()
	})
	executor.run()

	return c.markStatus(ClusterStatusRunning, multiErr.FinalError())
}

func (c *svcCluster) Stop() error {
	if status := c.Status(); status != ClusterStatusRunning {
		return errUnableToStopNotRunningCluster
	}

	var (
		lock     sync.Mutex
		multiErr xerrors.MultiError
	)

	executor := newConcurrentNodeExecutor(c.knownNodes, c.opts.NodeConcurrency(), func(node env.ServiceNode) {
		if err := node.Stop(); err != nil {
			lock.Lock()
			multiErr = multiErr.Add(err)
			lock.Unlock()
			return
		}
		lock.Lock()
		c.spares = append(c.spares, node)
		lock.Unlock()
	})
	executor.run()

	return c.markStatus(ClusterStatusInitialized, multiErr.FinalError())
}

func (c *svcCluster) Status() Status {
	c.statusLock.RLock()
	defer c.statusLock.RUnlock()
	return c.status
}
