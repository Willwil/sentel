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
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	swarm "github.com/docker/docker/client"
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

type swarmCluster struct {
	config         core.Config
	mutex          sync.Mutex
	brokers        map[string]*broker
	client         *swarm.Client
	requiredImages map[string]bool
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
		for _, name := range names {
			images[name] = false
		}
	}
	// Conenct with swarm
	cli, err := swarm.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("cluster manager failed to connect with swarm:'%s'", err.Error())
	}
	return &swarmCluster{
		config:         c,
		mutex:          sync.Mutex{},
		brokers:        make(map[string]*broker),
		requiredImages: images,
		client:         cli,
	}, nil

}

// Initiialize confirm and initialize cluster manager's prerequirment
func (p *swarmCluster) Initialize() error {
	// We must confirm wether the requrired images already exist
	ok := true
	for imageName, _ := range p.requiredImages {
		filters := filters.NewArgs()
		filters.Add("reference", imageName)
		options := types.ImageListOptions{
			Filters: filters,
		}
		_, err := p.client.ImageList(context.Background(), options)
		if err != nil {
			glog.Errorf("swarm cluster can not find required '%s docker image", imageName)
			ok = false
		}
	}
	if !ok {
		return errors.New("swarm cluster initialization failed to find required docker images")
	}
	return nil
}

// CreateBrokers create a number of brokers for tenant and product
func (p *swarmCluster) CreateBrokers(tid string, pid string, replicas int32) ([]string, error) {
	result := []string{}
	for i := 0; i < int(replicas); i++ {
		bid := fmt.Sprintf("%s-%s-%d", tid, pid, i)
		name := "sentel/broker"
		containerConfig := container.Config{}
		hostConfig := container.HostConfig{}
		netConfig := network.NetworkingConfig{}
		body, err := p.client.ContainerCreate(context.Background(), &containerConfig, &hostConfig, &netConfig, name)
		if err != nil {
			return nil, fmt.Errorf("swarm cluster failed to create docker container for '%s'", name)
		}
		result = append(result, name)
		p.mutex.Lock()
		p.brokers[bid] = &broker{
			bid:       bid,
			tid:       tid,
			pid:       pid,
			status:    brokerStatusCreated,
			createdAt: time.Now(),
			context:   body.ID,
		}
		p.mutex.Unlock()
	}
	return result, nil
}

// StartBroker start the specified broker
func (p *swarmCluster) StartBroker(bid string) error {
	if _, found := p.brokers[bid]; !found {
		return fmt.Errorf("Invalid broker id '%s' in swarm cluster", bid)
	}
	containerId := p.brokers[bid].context.(string)
	options := types.ContainerStartOptions{}
	if err := p.client.ContainerStart(context.Background(), containerId, options); err != nil {
		return fmt.Errorf("swarm failed to start broker '%s', reason:'%s'", bid, err.Error())
	}
	return nil
}

// StopBroker stop the specified broker
func (p *swarmCluster) StopBroker(bid string) error {
	if _, found := p.brokers[bid]; !found {
		return fmt.Errorf("Invalid broker id '%s' in swarm cluster", bid)
	}
	containerId := p.brokers[bid].context.(string)
	timeout := time.Duration(5 * time.Second)
	if err := p.client.ContainerStop(context.Background(), containerId, &timeout); err != nil {
		return fmt.Errorf("swarm failed to stop broker '%s', reason:'%s'", bid, err.Error())
	}
	return nil
}

// DeleteBroker stop and delete specified broker
func (p *swarmCluster) DeleteBroker(bid string) error {
	if _, found := p.brokers[bid]; !found {
		return fmt.Errorf("Invalid broker id '%s' in swarm cluster", bid)
	}
	containerId := p.brokers[bid].context.(string)
	options := types.ContainerRemoveOptions{
		Force: true,
	}
	if err := p.client.ContainerRemove(context.Background(), containerId, options); err != nil {
		return fmt.Errorf("swarm failed to remove broker '%s', reason:'%s'", bid, err.Error())
	}
	return nil

}

// RollbackBrokers rollback tenant's brokers
func (p *swarmCluster) RollbackBrokers(tid, pid string, replicas int32) error {
	return nil
}
