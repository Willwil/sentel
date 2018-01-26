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
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cloustone/sentel/pkg/config"
	sd "github.com/cloustone/sentel/pkg/service-discovery"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

type swarmCluster struct {
	config           config.Config
	mutex            sync.Mutex
	services         map[string]string
	client           *client.Client
	serviceDiscovery sd.ServiceDiscovery
}

func newSwarmCluster(c config.Config) (*swarmCluster, error) {
	cli, err := client.NewEnvClient()
	if err != nil || cli == nil {
		return nil, fmt.Errorf("cluster manager failed to connect with swarm:'%s'", err.Error())
	}
	return &swarmCluster{
		config:   c,
		mutex:    sync.Mutex{},
		services: make(map[string]string),
		client:   cli,
	}, nil

}

func (p *swarmCluster) SetServiceDiscovery(s sd.ServiceDiscovery) { p.serviceDiscovery = s }
func (p *swarmCluster) Initialize() error                         { return nil }

func (p *swarmCluster) CreateNetwork(name string) (string, error) {
	// Before create network, we should check wether the network already exist
	rsp, err := p.client.NetworkCreate(context.Background(),
		name,
		types.NetworkCreate{Driver: "overlay", CheckDuplicate: true})
	return rsp.ID, err
}

func (p *swarmCluster) RemoveNetwork(name string) error {
	return p.client.NetworkRemove(context.Background(), name)
}

func (p *swarmCluster) CreateService(spec ServiceSpec) (string, error) {
	delay := time.Duration(1 * time.Second)
	maxAttempts := uint64(10)
	reps := uint64(spec.Replicas)
	serviceName := spec.ServiceName

	service := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: serviceName,
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: spec.Image,
				Env:   spec.Environment,
			},
			Networks: []swarm.NetworkAttachmentConfig{
				{Target: spec.NetworkId},
			},
			RestartPolicy: &swarm.RestartPolicy{
				Condition:   swarm.RestartPolicyConditionOnFailure,
				Delay:       &delay,
				MaxAttempts: &maxAttempts,
			},
		},
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &reps,
			},
		},
	}
	options := types.ServiceCreateOptions{}
	serviceID := ""
	if rsp, err := p.client.ServiceCreate(context.Background(), service, options); err != nil {
		glog.Error(err)
		return serviceName, fmt.Errorf("swarm failed to create service '%s'", serviceName)
	} else {
		p.services[serviceName] = rsp.ID
		serviceID = rsp.ID
	}

	// update service discovery
	if p.serviceDiscovery != nil {
		service, _, err := p.client.ServiceInspectWithRaw(context.Background(), serviceID, types.ServiceInspectOptions{})
		if err != nil {
			return serviceName, fmt.Errorf("swarm failed to get service backend info for '%s'", serviceName)
		}

		if len(service.Spec.TaskTemplate.Networks) != 1 {
			return serviceName, fmt.Errorf("inspect service network failed")
		}
		networkID := service.Spec.TaskTemplate.Networks[0].Target

		ip := ""
		for _, vip := range service.Endpoint.VirtualIPs {
			if vip.NetworkID == networkID {
				ip = strings.Split(vip.Addr, "/")[0]
				break
			}
		}
		if ip == "" {
			return serviceName, fmt.Errorf("inspect service virtual ip failed")
		}

		if len(service.Endpoint.Ports) != 1 {
			return serviceName, fmt.Errorf("inspect service port failed")
		}

		dockerService := sd.Service{
			Name: serviceName,
			ID:   serviceID,
			IP:   ip,
			Port: service.Endpoint.Ports[0].PublishedPort,
		}
		if err := p.serviceDiscovery.RegisterService(dockerService); err != nil {
			return serviceName, fmt.Errorf("swarm failed to update service discovery for '%s'", serviceName)
		}
	}
	return serviceName, nil
}

func (p *swarmCluster) RemoveService(serviceName string) error {
	if id, found := p.services[serviceName]; found {
		if err := p.client.ServiceRemove(context.Background(), id); err != nil {
			return fmt.Errorf("swarm stop service '%s' failed", serviceName)
		}
		if p.serviceDiscovery != nil {
			p.serviceDiscovery.RemoveService(sd.Service{Name: serviceName, ID: id})
		}
	}
	return nil
}

func (p *swarmCluster) UpdateService(serviceId string, spec ServiceSpec) error {
	return nil
}

func (p *swarmCluster) IntrospectService(serviceID string) (ServiceIntrospec, error) {
	serviceSpec := ServiceIntrospec{
		ServiceId: serviceID,
	}
	if p.serviceDiscovery != nil {
		service, _, err := p.client.ServiceInspectWithRaw(context.Background(), serviceID, types.ServiceInspectOptions{})
		if err != nil {
			return serviceSpec, fmt.Errorf("swarm failed to get service backend info for '%s'", serviceID)
		}
		serviceSpec.ServiceName = service.Spec.Annotations.Name
		if len(service.Endpoint.Ports) == len(service.Endpoint.VirtualIPs) {
			for index, ep := range service.Endpoint.Ports {
				port := ep.PublishedPort
				vip := service.Endpoint.VirtualIPs[index].Addr
				serviceSpec.Endpoints = append(serviceSpec.Endpoints, ServiceEndpoint{VirtualIP: vip, Port: port})
			}
		}
	}
	return serviceSpec, nil
}
