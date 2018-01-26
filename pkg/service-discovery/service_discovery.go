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

package sd

import (
	"errors"

	"github.com/cloustone/sentel/pkg/config"
)

const (
	// BackendZookeeper zookeeper
	BackendZookeeper = "zookeeper"
	// BackendEtcd etcd
	BackendEtcd = "etcd"
)

// Service broker service define
type Service struct {
	Name string `json:"Name"`
	ID   string `json:"ID"`
	IP   string `json:"IP"`
	Port uint32 `json:"Port"`
}

// ServiceDiscovery service discovery
type ServiceDiscovery interface {
	RegisterService(s Service) error
	RemoveService(s Service)
	UpdateService(s Service) error
	Close()
}

// New create service discovery
func New(c config.Config, bks string) (ServiceDiscovery, error) {
	switch bks {
	case BackendZookeeper:
		return newServiceDiscoveryZK(c)
	default:
		return nil, errors.New("no valid service discovery backend")
	}
}
