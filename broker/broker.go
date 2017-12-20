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
	"errors"
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
	config   com.Config                // Global config
	services map[string]base.Service   // All service created by config.Protocols
	quits    map[string]chan os.Signal // Notification channel for each service
	name     string                    // Name of service manager
	brokerId string
}

const (
	BrokerVersion = "0.1"
)

var (
	serviceFactories = make(map[string]base.ServiceFactory)
	serviceSeqs      = []string{}
)

// RegisterService register service with name and protocol specified
func registerService(name string, factory base.ServiceFactory) {
	glog.Infof("Service '%s' is registered", name)
	if _, ok := serviceFactories[name]; ok {
		glog.Errorf("Service '%s' is already registered", name)
	}
	serviceFactories[name] = factory
	serviceSeqs = append(serviceSeqs, name)
}

// newBroker create global broker
func NewBroker(c com.Config) (*Broker, error) {
	broker := &Broker{
		config:   c,
		quits:    make(map[string]chan os.Signal),
		services: make(map[string]base.Service),
		brokerId: uuid.NewV4().String(),
	}
	base.SetBrokerStartupInfo(&base.BrokerStartupInfo{
		Id: broker.brokerId,
	})

	for name, _ := range serviceFactories {
		// Format service name
		quit := make(chan os.Signal)
		service, err := serviceFactories[name](c, quit)
		if err != nil {
			glog.Infof("Create service '%s'failed", name)
			return nil, err
		} else {
			glog.Infof("Create service '%s' successfully", name)
			broker.services[name] = service
			broker.quits[name] = quit
			base.RegisterService(name, service)
		}
	}
	return broker, nil
}

// getServicesByName return service instance by name, or matched by part of name
func (p *Broker) getServiceByName(name string) com.Service {
	if _, ok := p.services[name]; !ok {
		panic(fmt.Sprintf("Failed to find service '%s' in broker", name))
	}
	return p.services[name]
}

// Start
func (p *Broker) Run() error {
	// initialize event manager at first
	for _, name := range serviceSeqs {
		if err := p.services[name].Initialize(); err != nil {
			return err
		}
		glog.Infof("Initializing service '%s' ...successfuly", name)
	}

	// start each registered services
	for _, name := range serviceSeqs {
		if err := p.services[name].Start(); err != nil {
			return err
		}
		glog.Infof("Starting service '%s' ...successfuly", name)
	}
	// Wait all service to terminate in main context<TODO>
	for name, quit := range p.quits {
		<-quit
		glog.Info("Servide(%s) is terminated", name)
	}

	return nil
}

// Stop
func (p *Broker) Stop() {
	// Wait all service to terminate in main context<TODO>
	for name, quit := range p.quits {
		signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT)
		glog.Info("Servide(%s) is be terminating", name)
	}

}

// CheckAllRegisteredServices check all registered service simplily
func checkAllRegisteredServices() error {
	if len(serviceFactories) == 0 {
		return errors.New("No service registered")
	}
	for name, _ := range serviceFactories {
		glog.Infof("Service '%s' is registered", name)
	}
	return nil
}
