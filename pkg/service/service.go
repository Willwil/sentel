//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations

package service

import (
	"errors"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/golang/glog"
)

type Service interface {
	Name() string
	Initialize() error
	Start() error
	Stop()
}

type ServiceFactory interface {
	New(c config.Config) (Service, error)
}

type ServiceManager struct {
	sync.Once
	Config           config.Config // Global config
	serviceFactories []ServiceFactory
	services         []Service // All service created by config.Protocols
	name             string    // Name of service manager
	signalChan       chan os.Signal
	waitChan         chan interface{}
}

const serviceManagerVersion = "0.1"

var (
	serviceManager *ServiceManager
)

// GetServiceManager create service manager and all supported service
// The function should be called in service
func GetServiceManager() *ServiceManager { return serviceManager }

// NewServiceManager create ServiceManager only in main context
func NewServiceManager(name string, c config.Config) (*ServiceManager, error) {
	if serviceManager != nil {
		return serviceManager, errors.New("NewServiceManager had been called many times")
	}
	mgr := &ServiceManager{
		Config:           c,
		serviceFactories: []ServiceFactory{},
		services:         []Service{},
		name:             name,
		signalChan:       make(chan os.Signal),
		waitChan:         make(chan interface{}),
	}
	signal.Notify(mgr.signalChan, syscall.SIGINT, syscall.SIGTERM)
	serviceManager = mgr
	return serviceManager, nil
}

func signalHandler() {
	mgr := GetServiceManager()
	for {
		select {
		case <-mgr.signalChan:
			mgr.stop()
			mgr.waitChan <- "finished"
		}
	}
}

func (p *ServiceManager) AddService(factory ServiceFactory) *ServiceManager {
	p.serviceFactories = append(p.serviceFactories, factory)
	return p
}

// Run launch all serices and wait to terminate
func (p *ServiceManager) RunAndWait() error {
	// Create all services
	for _, factory := range p.serviceFactories {
		if service, err := factory.New(p.Config); err != nil {
			return err
		} else {
			glog.Infof("service '%s' created", service.Name())
			p.services = append(p.services, service)
		}
	}
	// Initialize all services
	for _, service := range p.services {
		if err := service.Initialize(); err != nil {
			glog.Infof("service '%s' initialize failed", service.Name())
			return err
		}
	}

	// Run all service
	glog.Infof("there are %d services in '%s'", len(p.services), p.name)
	for _, service := range p.services {
		if err := service.Start(); err != nil {
			return err
		}
		glog.Infof("service '%s' started", service.Name())
	}
	go signalHandler()
	<-p.waitChan
	return nil
}

// Run launch all serices
func (p *ServiceManager) Run() error {
	// Run all service
	glog.Infof("There are %d services in '%s'", len(p.services), p.name)
	for _, service := range p.services {
		glog.Infof("Starting service '%s'...", service.Name())
		if err := service.Start(); err != nil {
			return err
		}
	}
	return nil
}

// stop
func (p *ServiceManager) stop() error {
	for _, service := range p.services {
		service.Stop()
		glog.Infof("service '%s' stoped", service.Name())
	}
	return nil
}

// GetServicesByName return service instance by name, or matched by part of name
func (p *ServiceManager) GetServicesByName(name string) []Service {
	services := []Service{}
	for _, service := range p.services {
		if strings.IndexAny(service.Name(), name) >= 0 {
			services = append(services, service)
		}
	}
	return services
}

// GetService return service instance by name, or matched by part of name
func (p *ServiceManager) GetService(name string) Service {
	for _, service := range p.services {
		if service.Name() == name {
			return service
		}
	}
	return nil
}
