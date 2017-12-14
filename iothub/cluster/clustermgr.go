//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package cluster

import (
	"errors"

	"github.com/cloustone/sentel/core"
)

// ClusterManager is wrapper cluster manager built on top of swarm and kubernetes
type ClusterManager interface {
	// Initialize precheck wether cluster proconditions are meet
	Initialize() error
	// CreateNetwork create tenant work
	CreateNetwork(name string) (string, error)
	// RemoveNetwork remove tenant work
	RemoveNetwork(name string) error
	// CreateService create broker service for product
	CreateService(tid string, pid string, replicas int32) (string, error)
	// RemoveService remove broker service of product
	RemoveService(serviceName string) error
	// UpdateService updatet product' service about replicas and ...
	UpdateService(serviceName string, replicas int32) error
}

// New retrieve clustermanager instance connected with clustermgr
func New(c core.Config) (ClusterManager, error) {
	if v, err := c.String("iothub", "cluster"); err != nil {
		var cluster ClusterManager
		switch v {
		case "k8s":
			cluster, _ = newK8sCluster(c)
		case "swarm":
			cluster, _ = newSwarmCluster(c)
		}
		if cluster != nil {
			return cluster, cluster.Initialize()
		}
	}
	return nil, errors.New("iothub cluster manager initialize failed")
}
