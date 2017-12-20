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

package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/common"
	"github.com/golang/glog"
	uuid "github.com/satori/go.uuid"
)

type Broker struct {
	sync.Once
	config           com.Config            // Global config
	serviceFactories []base.ServiceFactory // All service factory
	services         []base.Service        // All service created by config.Protocols
	quits            []chan os.Signal      // Notification channel for each service
	name             string                // Name of service manager
	brokerId         string
}

const (
	BrokerVersion = "0.1"
)

// newBroker create global broker
func NewBroker(c com.Config) (*Broker, error) {
	broker := &Broker{
		config:           c,
		quits:            []chan os.Signal{},
		services:         []base.Service{},
		serviceFactories: []base.ServiceFactory{},
		brokerId:         uuid.NewV4().String(),
	}
	base.SetBrokerStartupInfo(&base.BrokerStartupInfo{
		Id: broker.brokerId,
	})

	return broker, nil
}

// addService register service with name and protocol specified
func (p *Broker) AddService(factory base.ServiceFactory) {
	p.serviceFactories = append(p.serviceFactories, factory)
}

// getServicesByName return service instance by name, or matched by part of name
func (p *Broker) GetService(name string) com.Service {
	for _, service := range p.services {
		if service.Name() == name {
			return service
		}
	}
	glog.Fatal(fmt.Errorf("Failed to find service '%s' in broker", name))
	return nil
}

// Start
func (p *Broker) Run() error {
	// create service
	for _, factory := range p.serviceFactories {
		// Format service name
		quit := make(chan os.Signal)
		service, err := factory.New(p.config, quit)
		if err != nil {
			return err
		} else {
			glog.Infof("Create service '%s' successfully", service.Name())
			p.services = append(p.services, service)
			p.quits = append(p.quits, quit)
			base.RegisterService(service.Name(), service)
		}
	}

	// initialize services
	for _, service := range p.services {
		if err := service.Initialize(); err != nil {
			return err
		}
		glog.Infof("Initializing service '%s' ...successfuly", service.Name())
	}

	// start each registered services
	for _, service := range p.services {
		if err := service.Start(); err != nil {
			return err
		}
		glog.Infof("Starting service '%s' ...successfuly", service.Name())
	}
	// Wait all service to terminate in main context<TODO>
	for index, service := range p.services {
		<-p.quits[index]
		glog.Info("Servide(%s) is terminated", service.Name())
	}

	return nil
}

// Stop
func (p *Broker) Stop() {
	// Wait all service to terminate in main context<TODO>
	for index, service := range p.services {
		signal.Notify(p.quits[index], syscall.SIGINT, syscall.SIGQUIT)
		glog.Info("Service(%s) is terminated", service.Name())
	}

}
