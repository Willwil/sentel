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

	"github.com/cloustone/sentel/pkg/config"
	sd "github.com/cloustone/sentel/pkg/service-discovery"
)

const (
	ServiceStateNone    = "none"
	ServiceStateStarted = "started"
	ServiceStateStoped  = "stoped"
)

type ServiceEndpoint struct {
	VirtualIP string
	Port      uint32
}

type ServiceSpec struct {
	TenantId    string
	ServiceName string
	Image       string
	NetworkId   string
	Replicas    int32
	Environment []string
}

type ServiceIntrospec struct {
	ServiceName  string
	ServiceId    string
	ServiceState string
	Endpoints    []ServiceEndpoint
}

// ClusterManager is wrapper cluster manager built on top of swarm and kubernetes
type ClusterManager interface {
	// SetDiscoveryBackend set service discovery backend
	SetServiceDiscovery(sd.ServiceDiscovery)
	// CreateNetwork create tenant network
	CreateNetwork(name string) (string, error)
	// RemoveNetwork remove tenant network
	RemoveNetwork(name string) error
	// CreateService create broker service for product
	CreateService(spec ServiceSpec) (string, error)
	// RemoveService remove broker service of product
	RemoveService(serviceId string) error
	// UpdateService updatet product' service about replicas and ...
	UpdateService(serviceId string, spec ServiceSpec) error
	// IntrospectService return service detail from cluster
	IntrospectService(serviceId string) (ServiceIntrospec, error)
}

// New retrieve clustermanager instance connected with clustermgr
func New(c config.Config) (ClusterManager, error) {
	if v, err := c.String("iothub", "cluster"); err == nil {
		switch v {
		case "k8s":
			return newK8sCluster(c)
		case "swarm":
			return newSwarmCluster(c)
		case "local":
			return newLocalCluster(c)
		}
	}
	return nil, errors.New("iothub cluster manager initialize failed")
}
