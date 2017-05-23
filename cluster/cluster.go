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

	"github.com/m3db/m3em/node"

	"github.com/m3db/m3cluster/services"
	"github.com/m3db/m3x/errors"
	"github.com/m3db/m3x/log"
	"github.com/m3db/m3x/sync"
)

var (
	errInsufficientCapacity          = fmt.Errorf("insufficient node capacity in environment")
	errNodeAlreadyUsed               = fmt.Errorf("unable to add node, already in use")
	errNodeNotInUse                  = fmt.Errorf("unable to remove node, not in use")
	errClusterNotUnitialized         = fmt.Errorf("unable to setup cluster, it is not unitialized")
	errClusterUnableToInitialize     = fmt.Errorf("unable to initialize cluster, it needs to be setup")
	errClusterUnableToAlterPlacement = fmt.Errorf("unable to alter cluster placement, it needs to be setup/running")
	errUnableToStartUnsetupCluster   = fmt.Errorf("unable to start cluster, it has not been setup")
	errClusterUnableToTeardown       = fmt.Errorf("unable to teardown cluster, it has not been setup")
	errUnableToStopNotRunningCluster = fmt.Errorf("unable to stop cluster, it is running")
)

type idToNodeMap map[string]node.ServiceNode

func (im idToNodeMap) values() []node.ServiceNode {
	returnNodes := make([]node.ServiceNode, 0, len(im))
	for _, node := range im {
		returnNodes = append(returnNodes, node)
	}
	return returnNodes
}

type svcCluster struct {
	sync.RWMutex

	logger       xlog.Logger
	opts         Options
	knownNodes   node.ServiceNodes
	usedNodes    idToNodeMap
	spares       []node.ServiceNode
	sparesByID   map[string]node.ServiceNode
	placementSvc services.PlacementService
	placement    services.ServicePlacement
	status       Status
	lastErr      error
}

// New returns a new cluster backed by provided service nodes
func New(
	nodes node.ServiceNodes,
	opts Options,
) (Cluster, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	cluster := &svcCluster{
		logger:       opts.InstrumentOptions().Logger(),
		opts:         opts,
		knownNodes:   nodes,
		usedNodes:    make(idToNodeMap, len(nodes)),
		spares:       make([]node.ServiceNode, 0, len(nodes)),
		sparesByID:   make(map[string]node.ServiceNode, len(nodes)),
		placementSvc: opts.PlacementService(),
		status:       ClusterStatusUninitialized,
	}
	cluster.addSparesWithLock(nodes)

	return cluster, nil
}

func (c *svcCluster) addSparesWithLock(spares []node.ServiceNode) {
	for _, spare := range spares {
		c.spares = append(c.spares, spare)
		c.sparesByID[spare.ID()] = spare
	}
}

func nodeSliceWithoutID(originalSlice node.ServiceNodes, removeID string) node.ServiceNodes {
	newSlice := make(node.ServiceNodes, 0, len(originalSlice))
	for _, elem := range originalSlice {
		if elem.ID() != removeID {
			newSlice = append(newSlice, elem)
		}
	}
	return newSlice
}

// TODO(prateek): reset initial placement after teardown
// TODO(prateek): use concurrency and other options

type concurrentNodeFn func(node.ServiceNode)

type concurrentNodeExecutor struct {
	wg      sync.WaitGroup
	nodes   node.ServiceNodes
	workers xsync.WorkerPool
	fn      concurrentNodeFn
}

func newConcurrentNodeExecutor(
	nodes node.ServiceNodes,
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

func (c *svcCluster) Placement() services.ServicePlacement {
	c.Lock()
	defer c.Unlock()
	return c.placement
}

func (c *svcCluster) initWithLock() error {
	psvc := c.placementSvc
	_, _, err := psvc.Placement()
	if err != nil { // attempt to retrieve current placement
		c.logger.Infof("unable to retrieve existing placement, skipping delete attempt")
	} else {
		// delete existing placement
		err = psvc.Delete()
		if err != nil {
			return fmt.Errorf("unable to delete existing placement during setup(): %+v", err)
		}
		c.logger.Infof("successfully deleted existing placement")
	}

	var (
		svcBuild        = c.opts.ServiceBuild()
		svcConf         = c.opts.ServiceConfig()
		sessionToken    = c.opts.SessionToken()
		sessionOverride = c.opts.SessionOverride()
		listener        = c.opts.NodeListener()
		lock            sync.Mutex
		multiErr        xerrors.MultiError
	)

	// setup all known service nodes with build, config
	executor := newConcurrentNodeExecutor(c.knownNodes, c.opts.NodeConcurrency(), func(node node.ServiceNode) {
		err := node.Setup(svcBuild, svcConf, sessionToken, sessionOverride)
		if err != nil {
			lock.Lock()
			multiErr = multiErr.Add(err)
			lock.Unlock()
			return
		}
		if listener != nil {
			node.RegisterListener(listener)
			// TODO(prateek): track listenerID returned in cluster struct, cleanup in Teardown()
		}
	})
	executor.run()

	return multiErr.FinalError()
}

func (c *svcCluster) Setup(numNodes int) ([]node.ServiceNode, error) {
	c.Lock()
	defer c.Unlock()

	if c.status != ClusterStatusUninitialized {
		return nil, errClusterNotUnitialized
	}

	if err := c.initWithLock(); err != nil {
		return nil, err
	}

	numSpares := len(c.spares)
	if numSpares < numNodes {
		return nil, errInsufficientCapacity
	}

	psvc := c.placementSvc
	placement, err := psvc.BuildInitialPlacement(c.sparesAsPlacementInstaceWithLock(), c.opts.NumShards(), c.opts.Replication())
	if err != nil {
		return nil, err
	}

	// update ServiceNode with new shards from placement
	var (
		multiErr      xerrors.MultiError
		usedInstances = placement.Instances()
		setupNodes    = make([]node.ServiceNode, 0, len(usedInstances))
	)
	for _, instance := range usedInstances {
		setupNode, err := c.markSpareUsedWithLock(instance)
		if err != nil {
			multiErr = multiErr.Add(err)
			continue
		}
		setupNodes = append(setupNodes, setupNode)
	}

	multiErr = multiErr.
		Add(c.setPlacementWithLock(placement))

	return setupNodes, c.markStatusWithLock(ClusterStatusSetup, multiErr.FinalError())
}

func (c *svcCluster) markSpareUsedWithLock(spare services.PlacementInstance) (node.ServiceNode, error) {
	id := spare.ID()
	spareNode, ok := c.sparesByID[id]
	if !ok {
		// should never happen
		return nil, fmt.Errorf("unable to find spare node with id: %s", id)
	}
	delete(c.sparesByID, id)
	c.spares = nodeSliceWithoutID(c.spares, id)
	c.usedNodes[id] = spareNode
	return spareNode, nil
}

func (c *svcCluster) AddNode() (node.ServiceNode, error) {
	c.Lock()
	defer c.Unlock()

	if c.status != ClusterStatusRunning && c.status != ClusterStatusSetup {
		return nil, errClusterUnableToAlterPlacement
	}

	numSpares := len(c.spares)
	if numSpares < 1 {
		return nil, errInsufficientCapacity
	}

	psvc := c.placementSvc
	newPlacement, usedInstance, err := psvc.AddInstance(c.sparesAsPlacementInstaceWithLock())
	if err != nil {
		return nil, err
	}

	setupNode, err := c.markSpareUsedWithLock(usedInstance)
	if err != nil {
		return nil, err
	}

	c.setPlacementWithLock(newPlacement)
	return setupNode, nil
}

func (c *svcCluster) setPlacementWithLock(p services.ServicePlacement) error {
	for _, instance := range p.Instances() {
		// nb(prateek): update usedNodes with the new shards.
		instanceID := instance.ID()
		usedNode, ok := c.usedNodes[instanceID]
		if ok {
			usedNode.SetShards(instance.Shards())
		}
	}

	c.placement = p
	return nil
}

func (c *svcCluster) sparesAsPlacementInstaceWithLock() []services.PlacementInstance {
	spares := make([]services.PlacementInstance, 0, len(c.spares))
	for _, spare := range c.spares {
		spares = append(spares, spare.(services.PlacementInstance))
	}
	return spares
}

func (c *svcCluster) RemoveNode(i node.ServiceNode) error {
	c.Lock()
	defer c.Unlock()

	if c.status != ClusterStatusRunning && c.status != ClusterStatusSetup {
		return errClusterUnableToAlterPlacement
	}

	usedNode, ok := c.usedNodes[i.ID()]
	if !ok {
		return errNodeNotInUse
	}

	psvc := c.placementSvc
	newPlacement, err := psvc.RemoveInstance(i.ID())
	if err != nil {
		return err
	}

	// update removed instance from used -> spare
	// nb(prateek): this omits modeling "leaving" shards on the node being removed
	usedNode.SetShards(nil)
	delete(c.usedNodes, usedNode.ID())
	c.addSparesWithLock([]node.ServiceNode{usedNode})

	return c.setPlacementWithLock(newPlacement)
}

func (c *svcCluster) ReplaceNode(oldNode node.ServiceNode) ([]node.ServiceNode, error) {
	c.Lock()
	defer c.Unlock()

	if c.status != ClusterStatusRunning && c.status != ClusterStatusSetup {
		return nil, errClusterUnableToAlterPlacement
	}

	if _, ok := c.usedNodes[oldNode.ID()]; !ok {
		return nil, errNodeNotInUse
	}

	psvc := c.placementSvc
	newPlacement, newInstances, err := psvc.ReplaceInstance(oldNode.ID(), c.sparesAsPlacementInstaceWithLock())
	if err != nil {
		return nil, err
	}

	// mark old node no longer used
	oldNode.SetShards(nil)
	delete(c.usedNodes, oldNode.ID())
	c.addSparesWithLock([]node.ServiceNode{oldNode})

	var (
		multiErr xerrors.MultiError
		newNodes = make([]node.ServiceNode, 0, len(newInstances))
	)
	for _, instance := range newInstances {
		newNode, err := c.markSpareUsedWithLock(instance)
		if err != nil {
			multiErr = multiErr.Add(err)
			continue
		}
		newNodes = append(newNodes, newNode)
	}

	multiErr = multiErr.
		Add(c.setPlacementWithLock(newPlacement))

	return newNodes, nil
}

func (c *svcCluster) SpareNodes() []node.ServiceNode {
	c.Lock()
	defer c.Unlock()
	return c.spares
}

func (c *svcCluster) ActiveNodes() []node.ServiceNode {
	c.Lock()
	defer c.Unlock()
	return c.usedNodes.values()
}

func (c *svcCluster) KnownNodes() []node.ServiceNode {
	c.Lock()
	defer c.Unlock()
	return c.knownNodes
}

func (c *svcCluster) markStatusWithLock(status Status, err error) error {
	if err == nil {
		c.status = status
		return nil
	}

	c.status = ClusterStatusError
	c.lastErr = err
	return err
}

func (c *svcCluster) Teardown() error {
	c.Lock()
	defer c.Unlock()

	if c.status == ClusterStatusUninitialized {
		return errClusterUnableToTeardown
	}

	var (
		lock     sync.Mutex
		multiErr xerrors.MultiError
	)

	executor := newConcurrentNodeExecutor(c.knownNodes, c.opts.NodeConcurrency(), func(node node.ServiceNode) {
		if err := node.Teardown(); err != nil {
			lock.Lock()
			multiErr = multiErr.Add(err)
			lock.Unlock()
			return
		}
	})
	executor.run()

	for id, usedNode := range c.usedNodes {
		usedNode.SetShards(nil)
		delete(c.usedNodes, id)
	}
	c.spares = make([]node.ServiceNode, 0, len(c.knownNodes))
	c.sparesByID = make(map[string]node.ServiceNode, len(c.knownNodes))
	c.addSparesWithLock(c.knownNodes)

	return c.markStatusWithLock(ClusterStatusUninitialized, multiErr.FinalError())
}

func (c *svcCluster) Start() error {
	c.Lock()
	defer c.Unlock()

	if c.status != ClusterStatusSetup {
		return errUnableToStartUnsetupCluster
	}

	var (
		lock     sync.Mutex
		multiErr xerrors.MultiError
	)

	executor := newConcurrentNodeExecutor(c.usedNodes.values(), c.opts.NodeConcurrency(), func(node node.ServiceNode) {
		if err := node.Start(); err != nil {
			lock.Lock()
			multiErr = multiErr.Add(err)
			lock.Unlock()
			return
		}
	})
	executor.run()

	return c.markStatusWithLock(ClusterStatusRunning, multiErr.FinalError())
}

func (c *svcCluster) Stop() error {
	c.Lock()
	defer c.Unlock()

	if c.status != ClusterStatusRunning {
		return errUnableToStopNotRunningCluster
	}

	var (
		lock     sync.Mutex
		multiErr xerrors.MultiError
	)

	executor := newConcurrentNodeExecutor(c.usedNodes.values(), c.opts.NodeConcurrency(), func(node node.ServiceNode) {
		if err := node.Stop(); err != nil {
			lock.Lock()
			multiErr = multiErr.Add(err)
			lock.Unlock()
			return
		}
	})
	executor.run()

	return c.markStatusWithLock(ClusterStatusSetup, multiErr.FinalError())
}

func (c *svcCluster) Status() Status {
	c.RLock()
	defer c.RUnlock()
	return c.status
}
