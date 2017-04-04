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

	env "github.com/m3db/m3em/environment"
	"github.com/m3db/m3cluster/services"
	m3dbclient "github.com/m3db/m3db/client"
	"github.com/m3db/m3x/errors"
	"github.com/m3db/m3x/log"
)

var (
	errInsufficientCapacity            = fmt.Errorf("insufficient instance capacity in environment")
	errInstanceAlreadyUsed             = fmt.Errorf("unable to add instance, already in use")
	errInstanceNotInUse                = fmt.Errorf("unable to remove instance, not in use")
	errClusterAlreadySetup             = fmt.Errorf("cluster already setup")
	errClusterUnableToInitialize       = fmt.Errorf("unable to initialize cluster, it needs to be setup")
	errClusterUnableToAlterPlacement   = fmt.Errorf("unable to alter cluster placement, it needs to be initialized/running")
	errClusterUnableToReset            = fmt.Errorf("unable to reset cluster, it needs to be setup|initialized|error")
	errUnableToStartUnitializedCluster = fmt.Errorf("unable to start unitialized cluster")
	errClusterUnableToTeardown         = fmt.Errorf("unable to teardown cluster, it has not been setup")
	errUnableToStopNotRunningCluster   = fmt.Errorf("unable to stop cluster, it is running")
)

type idToInstanceMap map[string]env.M3DBInstance

type m3dbCluster struct {
	logger        xlog.Logger
	copts         Options
	env           env.M3DBEnvironment
	usedInstances idToInstanceMap
	spares        []env.M3DBInstance
	placementSvc  services.PlacementService

	statusLock sync.RWMutex
	status     Status
	lastErr    error
}

// New returns a new M3DB cluster of instances backed by the provided environment
func New(m3env env.M3DBEnvironment, cOpts Options) Cluster {
	return &m3dbCluster{
		logger:        cOpts.InstrumentOptions().Logger(),
		copts:         cOpts,
		env:           m3env,
		usedInstances: make(idToInstanceMap, len(m3env.Instances())),
		placementSvc:  cOpts.PlacementService(),
		status:        ClusterStatusUninitialized,
	}
}

// TODO(prateek): reset initial placement after teardown
// TODO(prateek): use concurrency and other options

func (c *m3dbCluster) Setup() error {
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
		svcBuild  = c.copts.ServiceBuild()
		svcConf   = c.copts.ServiceConfig()
		wg        sync.WaitGroup
		lock      sync.Mutex
		instances = c.env.Instances()
		multiErr  xerrors.MultiError
	)

	// setup the instances, in parallel
	for idx := range instances {
		inst := instances[idx]
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := inst.Setup(svcBuild, svcConf); err != nil {
				lock.Lock()
				multiErr = multiErr.Add(err)
				lock.Unlock()
				return
			}

			lock.Lock()
			c.spares = append(c.spares, inst)
			lock.Unlock()
		}()
	}
	wg.Wait()

	return c.markStatus(ClusterStatusSetup, multiErr.FinalError())
}

func (c *m3dbCluster) Initialize(numNodes int) ([]env.M3DBInstance, error) {
	if c.Status() != ClusterStatusSetup {
		return nil, errClusterUnableToInitialize
	}

	numSpares := len(c.spares)
	if numSpares < numNodes {
		return nil, errInsufficientCapacity
	}

	usedInstances := c.spares[:numNodes]
	instances := make([]services.PlacementInstance, 0, numNodes)
	for _, inst := range usedInstances {
		c.usedInstances[inst.ID()] = inst
		instances = append(instances, inst)
	}
	c.spares = c.spares[numNodes:]

	psvc := c.placementSvc
	_, err := psvc.BuildInitialPlacement(instances, c.copts.NumShards(), c.copts.Replication())
	if err != nil {
		return nil, err
	}

	c.statusLock.Lock()
	c.status = ClusterStatusInitialized
	c.statusLock.Unlock()
	return usedInstances, nil
}

func (c *m3dbCluster) AddInstance() (env.M3DBInstance, error) {
	if status := c.Status(); status != ClusterStatusRunning && status != ClusterStatusInitialized {
		return nil, errClusterUnableToAlterPlacement
	}

	numSpares := len(c.spares)
	if numSpares < 1 {
		return nil, errInsufficientCapacity
	}
	inst := c.spares[0]
	psvc := c.placementSvc
	_, err := psvc.AddInstance([]services.PlacementInstance{inst})
	if err != nil {
		return nil, err
	}
	c.usedInstances[inst.ID()] = inst
	c.spares = c.spares[1:]
	return inst, nil
}

func (c *m3dbCluster) RemoveInstance(i env.M3DBInstance) error {
	if status := c.Status(); status != ClusterStatusRunning && status != ClusterStatusInitialized {
		return errClusterUnableToAlterPlacement
	}

	if _, ok := c.usedInstances[i.ID()]; !ok {
		return errInstanceNotInUse
	}
	psvc := c.placementSvc
	_, err := psvc.RemoveInstance(i.ID())
	if err != nil {
		return err
	}
	delete(c.usedInstances, i.ID())
	c.spares = append(c.spares, i)
	return nil
}

func (c *m3dbCluster) ReplaceInstance(oldInstance env.M3DBInstance) (env.M3DBInstance, error) {
	if status := c.Status(); status != ClusterStatusRunning && status != ClusterStatusInitialized {
		return nil, errClusterUnableToAlterPlacement
	}

	if _, ok := c.usedInstances[oldInstance.ID()]; !ok {
		return nil, errInstanceNotInUse
	}

	numSpares := len(c.spares)
	if numSpares < 1 {
		return nil, errInsufficientCapacity
	}

	spareInstance := c.spares[0]
	psvc := c.placementSvc
	_, err := psvc.ReplaceInstance(oldInstance.ID(), []services.PlacementInstance{spareInstance})
	if err != nil {
		return nil, err
	}

	delete(c.usedInstances, oldInstance.ID())
	c.usedInstances[spareInstance.ID()] = spareInstance
	c.spares = append(c.spares[1:], oldInstance)

	return spareInstance, nil
}

func (c *m3dbCluster) Spares() []env.M3DBInstance {
	return c.spares
}

func (c *m3dbCluster) Reset() error {
	if status := c.Status(); status != ClusterStatusRunning && status != ClusterStatusInitialized &&
		status != ClusterStatusSetup && status != ClusterStatusError {
		return errClusterUnableToReset
	}

	var (
		wg        sync.WaitGroup
		lock      sync.Mutex
		instances = c.env.Instances()
		multiErr  xerrors.MultiError
	)

	// setup the instances, in parallel
	for idx := range instances {
		inst := instances[idx]
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := inst.Reset(); err != nil {
				lock.Lock()
				multiErr = multiErr.Add(err)
				lock.Unlock()
				return
			}

			lock.Lock()
			c.spares = append(c.spares, inst)
			lock.Unlock()
		}()
	}
	wg.Wait()

	return c.markStatus(ClusterStatusSetup, multiErr.FinalError())
}

func (c *m3dbCluster) markStatus(status Status, err error) error {
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

func (c *m3dbCluster) Teardown() error {
	if status := c.Status(); status == ClusterStatusUninitialized {
		return errClusterUnableToTeardown
	}

	var (
		wg        sync.WaitGroup
		lock      sync.Mutex
		instances = c.env.Instances()
		multiErr  xerrors.MultiError
	)

	// setup the instances, in parallel
	for idx := range instances {
		inst := instances[idx]
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := inst.Teardown(); err != nil {
				lock.Lock()
				multiErr = multiErr.Add(err)
				lock.Unlock()
				return
			}

			lock.Lock()
			c.spares = append(c.spares, inst)
			lock.Unlock()
		}()
	}
	wg.Wait()

	return c.markStatus(ClusterStatusUninitialized, multiErr.FinalError())
}

func (c *m3dbCluster) StartInitialized() error {
	if status := c.Status(); status != ClusterStatusInitialized {
		return errUnableToStartUnitializedCluster
	}

	var (
		wg        sync.WaitGroup
		lock      sync.Mutex
		instances = c.usedInstances
		multiErr  xerrors.MultiError
	)

	// setup the instances, in parallel
	for instanceID := range instances {
		inst := instances[instanceID]
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := inst.Start(); err != nil {
				lock.Lock()
				multiErr = multiErr.Add(err)
				lock.Unlock()
				return
			}

			lock.Lock()
			c.spares = append(c.spares, inst)
			lock.Unlock()
		}()
	}
	wg.Wait()

	return c.markStatus(ClusterStatusRunning, multiErr.FinalError())
}

func (c *m3dbCluster) Start() error {
	if status := c.Status(); status != ClusterStatusInitialized {
		return errUnableToStartUnitializedCluster
	}

	var (
		wg        sync.WaitGroup
		lock      sync.Mutex
		instances = c.env.Instances()
		multiErr  xerrors.MultiError
	)

	// setup the instances, in parallel
	for instanceID := range instances {
		inst := instances[instanceID]
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := inst.Start(); err != nil {
				lock.Lock()
				multiErr = multiErr.Add(err)
				lock.Unlock()
				return
			}

			lock.Lock()
			c.spares = append(c.spares, inst)
			lock.Unlock()
		}()
	}
	wg.Wait()

	return c.markStatus(ClusterStatusRunning, multiErr.FinalError())
}

func (c *m3dbCluster) Stop() error {
	if status := c.Status(); status != ClusterStatusRunning {
		return errUnableToStopNotRunningCluster
	}

	var (
		wg        sync.WaitGroup
		lock      sync.Mutex
		instances = c.env.Instances()
		multiErr  xerrors.MultiError
	)

	// setup the instances, in parallel
	for instanceID := range instances {
		inst := instances[instanceID]
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := inst.Stop(); err != nil {
				lock.Lock()
				multiErr = multiErr.Add(err)
				lock.Unlock()
				return
			}

			lock.Lock()
			c.spares = append(c.spares, inst)
			lock.Unlock()
		}()
	}
	wg.Wait()

	return c.markStatus(ClusterStatusInitialized, multiErr.FinalError())
}

func (c *m3dbCluster) Status() Status {
	c.statusLock.RLock()
	defer c.statusLock.RUnlock()
	return c.status
}

func (c *m3dbCluster) Client() m3dbclient.Client {
	panic("not implemented")
}

func (c *m3dbCluster) AdminClient() m3dbclient.AdminClient {
	panic("not implemented")
}