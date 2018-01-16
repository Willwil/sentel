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
	"fmt"
	"strings"
	"sync"
	"time"

	sd "github.com/cloustone/sentel/iothub/service-discovery"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/docker-service"
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
	// Get requried images
	images := make(map[string]bool)
	if v, err := c.String("iothub", "docker-images"); err != nil || v == "" {
		return nil, errors.New("invalid configuration for docker-images in iothub.conf")
	} else {
		names := strings.Split(v, ",")
		if len(names) == 0 {
			return nil, errors.New("no docker-images are specified in iothub.conf")
		}
		for _, v := range names {
			name := strings.TrimSpace(v)
			images[name] = false
		}
	}
	// Conenct with swarm
	cli, err := client.NewEnvClient()
	if err != nil || cli == nil {
		return nil, fmt.Errorf("cluster manager failed to connect with swarm:'%s'", err.Error())
	}
	/*
		for imageName, _ := range images {
			filters := filters.NewArgs()
			filters.Add("reference", imageName)
			options := types.ImageListOptions{
				Filters: filters,
			}
			if _, err := cli.ImageList(context.Background(), options); err != nil {
				return nil, fmt.Errorf("swarm cluster can not find required '%s docker image", imageName)
			}
			glog.Infof("swarm found docker image '%s' in docker host", imageName)
		}
	*/
	return &swarmCluster{
		config:   c,
		mutex:    sync.Mutex{},
		services: make(map[string]string),
		client:   cli,
	}, nil

}

func (p *swarmCluster) SetServiceDiscovery(s sd.ServiceDiscovery) {
	p.serviceDiscovery = s
}

func (p *swarmCluster) Initialize() error {
	/*
		// leave first
		//p.client.SwarmLeave(context.Background(), true)
		// become a swarm manager
		options := swarm.InitRequest{
			ListenAddr: "0.0.0.0:2377",
		}
		_, err := p.client.SwarmInit(context.Background(), options)
		return err
	*/
	return nil
}

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

func (p *swarmCluster) CreateService(tid string, replicas int32) (string, error) {
	serviceName := fmt.Sprintf("tenant_%s", tid)
	env := []string{
		"KAFKA_HOST=kafka:9092",
		"MONGO_HOST=mongo:27017",
		fmt.Sprintf("BROKER_TENANT=%s", tid),
	}
	delay := time.Duration(1 * time.Second)
	maxAttempts := uint64(10)
	reps := uint64(replicas)

	service := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: serviceName,
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: "sentel/broker",
				Env:   env,
			},
			Networks: []swarm.NetworkAttachmentConfig{
				{Target: tid},
				{Target: "sentel_front"},
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
		if len(service.Endpoint.Ports) == len(service.Endpoint.VirtualIPs) {
			// endpoints := []ds.ServiceEndpoint{}
			// for index, ep := range service.Endpoint.Ports {
			// 	port := ep.PublishedPort
			// 	vip := service.Endpoint.VirtualIPs[index].Addr
			// 	endpoints = append(endpoints, sd.ServiceEndpoint{IP: vip, Port: port})
			// }

			service := ds.Service{Name: serviceName, ID: serviceID}
			if err := p.serviceDiscovery.RegisterService(service); err != nil {
				return serviceName, fmt.Errorf("swarm failed to update service discovery for '%s'", serviceName)
			}
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
			p.serviceDiscovery.RemoveService(ds.Service{Name: serviceName, ID: id})
		}
	}
	return nil
}

func (p *swarmCluster) UpdateService(serviceName string, replicas int32) error {
	return nil
}

func (p *swarmCluster) IntrospectService(serviceID string) (ServiceSpec, error) {
	serviceSpec := ServiceSpec{
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
