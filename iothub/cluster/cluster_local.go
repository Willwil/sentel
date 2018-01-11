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
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"sync"

	sd "github.com/cloustone/sentel/iothub/service-discovery"
	"github.com/cloustone/sentel/pkg/config"
)

type localCluster struct {
	config           config.Config
	mutex            sync.Mutex
	services         map[string]*exec.Cmd
	serviceDiscovery sd.ServiceDiscovery
	ports            map[uint32]string
	portIndex        uint32
	serviceSpecs     map[string]ServiceSpec
	ctxs             map[string]context.Context
}

func newLocalCluster(c config.Config) (*localCluster, error) {
	return &localCluster{
		config:       c,
		mutex:        sync.Mutex{},
		services:     make(map[string]*exec.Cmd),
		ports:        make(map[uint32]string),
		portIndex:    10000,
		serviceSpecs: make(map[string]ServiceSpec),
		ctxs:         make(map[string]context.Context),
	}, nil
}

func (p *localCluster) makePort() uint32 {
	scope := 1000
	for i := 0; i < scope; i++ {
		port := p.portIndex + uint32(rand.Intn(scope))
		if _, found := p.ports[port]; !found {
			return port
		}
	}
	return 0
}

func (p *localCluster) SetServiceDiscovery(s sd.ServiceDiscovery) { p.serviceDiscovery = s }
func (p *localCluster) Initialize() error                         { return nil }
func (p *localCluster) CreateNetwork(name string) (string, error) { return "", nil }
func (p *localCluster) RemoveNetwork(name string) error           { return nil }

func (p *localCluster) CreateService(tid string, replicas int32) (string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// Now only support one service instance in local cluster mode
	port := p.makePort()
	ctx := context.Background()
	cmd := exec.CommandContext(ctx,
		"broker",
		"-d",
		fmt.Sprintf("-t %s", tid),
		"-P tcp",
		fmt.Sprintf("-l localhost:%d", port))
	spec := ServiceSpec{
		ServiceName:  tid,
		ServiceId:    fmt.Sprintf("%s", len(p.services)+1),
		ServiceState: ServiceStateStarted,
		Endpoints:    []ServiceEndpoint{{VirtualIP: "127.0.0.1", Port: uint32(port)}},
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}
	serviceId := spec.ServiceId
	p.services[serviceId] = cmd
	p.ports[port] = serviceId
	p.serviceSpecs[serviceId] = spec
	p.ctxs[serviceId] = ctx

	if p.serviceDiscovery != nil {
		service := sd.Service{
			ServiceName: tid,
			ServiceId:   serviceId,
			Endpoints:   []sd.ServiceEndpoint{{IP: "127.0.0.1", Port: port}},
		}
		p.serviceDiscovery.RegisterService(service)
	}
	return serviceId, nil
}

func (p *localCluster) RemoveService(serviceId string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, found := p.services[serviceId]; !found {
		return fmt.Errorf("service '%s' not found", serviceId)
	}
	ctx := p.ctxs[serviceId]
	ctx.Done()
	delete(p.services, serviceId)
	delete(p.ctxs, serviceId)
	delete(p.serviceSpecs, serviceId)
	spec := p.serviceSpecs[serviceId]
	delete(p.ports, spec.Endpoints[0].Port)
	if p.serviceDiscovery != nil {
		p.serviceDiscovery.RemoveService(sd.Service{ServiceId: serviceId})
	}
	return nil
}

func (p *localCluster) UpdateService(serviceId string, replicas int32) error {
	return nil
}

func (p *localCluster) IntrospectService(serviceId string) (ServiceSpec, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	serviceSpec := ServiceSpec{
		ServiceId: serviceId,
	}
	if _, found := p.serviceSpecs[serviceId]; !found {
		return serviceSpec, fmt.Errorf("no service '%s'", serviceId)
	}
	return p.serviceSpecs[serviceId], nil
}
