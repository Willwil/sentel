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

// ClusterManager manage service in cluster for tenant and product
type ClusterManager interface {
	Initialize() error
	CreateService(tid string, pid string, replicas int32) (string, error)
	RemoveService(serviceName string) error
	UpdateService(serviceName string, replicas int32) error
}

// New retrieve clustermanager instance connected with clustermgr
func New(c core.Config) (ClusterManager, error) {
	if v, err := c.String("iothub", "cluster"); err != nil {
		return nil, errors.New("iothub cluster manager is not specified")
	} else {
		var cluster ClusterManager
		var err error
		switch v {
		case "k8s":
			cluster, err = newK8sCluster(c)
		case "swarm":
			cluster, err = newSwarmCluster(c)
		}
		if err != nil {
			return nil, errors.New("iothub cluster manager is not specified")
		} else {
			return cluster, cluster.Initialize()
		}
	}
}
