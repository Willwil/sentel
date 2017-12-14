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

	"github.com/cloustone/sentel/core"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

type swarmCluster struct {
	config   core.Config
	mutex    sync.Mutex
	services map[string]string
	client   *client.Client
}

func newSwarmCluster(c core.Config) (*swarmCluster, error) {
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
	if err != nil {
		return nil, fmt.Errorf("cluster manager failed to connect with swarm:'%s'", err.Error())
	}
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
	return &swarmCluster{
		config:   c,
		mutex:    sync.Mutex{},
		services: make(map[string]string),
		client:   cli,
	}, nil

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
	rsp, err := p.client.NetworkCreate(context.Background(), name, types.NetworkCreate{Driver: "overlay"})
	return rsp.ID, err
}

func (p *swarmCluster) RemoveNetwork(name string) error {
	return p.client.NetworkRemove(context.Background(), name)
}

func (p *swarmCluster) CreateService(tid string, pid string, replicas int32) (string, error) {
	serviceName := fmt.Sprintf("%s-%s", pid, tid)
	env := []string{
		"KAFKA_HOST=sentel_kafka",
		"MONGO_HOST=sentel_mongo",
		fmt.Sprintf("BROKER_TENANT=%s", tid),
		fmt.Sprintf("BROKER_PRODUCT=%s", pid),
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
			//TODO: how to add network setting
			Networks: []swarm.NetworkAttachmentConfig{
				{Target: tid},
				{Target: "sentel-front"},
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
	if rsp, err := p.client.ServiceCreate(context.Background(), service, options); err != nil {
		glog.Error(err)
		return serviceName, fmt.Errorf("swarm failed to create service '%s'", serviceName)
	} else {
		p.services[serviceName] = rsp.ID
	}
	return serviceName, nil
}

func (p *swarmCluster) RemoveService(serviceName string) error {
	if id, found := p.services[serviceName]; found {
		if err := p.client.ServiceRemove(context.Background(), id); err != nil {
			return fmt.Errorf("swarm stop service '%s' failed", serviceName)
		}
	}
	return nil
}

func (p *swarmCluster) UpdateService(serviceName string, replicas int32) error {
	return nil
}
