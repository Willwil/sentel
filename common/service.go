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

package com

import (
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/golang/glog"
)

type Service interface {
	Name() string
	Start() error
	Stop()
}

type ServiceFactory interface {
	New(c Config, quit chan os.Signal) (Service, error)
}

type ServiceManager struct {
	sync.Once
	Config           Config // Global config
	serviceFactories map[string]ServiceFactory
	Services         map[string]Service        // All service created by config.Protocols
	quits            map[string]chan os.Signal // Notification channel for each service
	name             string                    // Name of service manager
}

const serviceManagerVersion = "0.1"

var (
	serviceManager *ServiceManager
)

// GetServiceManager create service manager and all supported service
// The function should be called in service
func GetServiceManager() *ServiceManager { return serviceManager }

// NewServiceManager create ServiceManager only in main context
func NewServiceManager(name string, c Config) (*ServiceManager, error) {
	if serviceManager != nil {
		return serviceManager, errors.New("NewServiceManager had been called many times")
	}
	mgr := &ServiceManager{
		Config:           c,
		serviceFactories: make(map[string]ServiceFactory),
		quits:            make(map[string]chan os.Signal),
		Services:         make(map[string]Service),
		name:             name,
	}
	serviceManager = mgr
	return serviceManager, nil
}

func (p *ServiceManager) AddService(name string, factory ServiceFactory) error {
	if _, found := p.serviceFactories[name]; found {
		glog.Fatalf("Service '%s' is already registered", name)
	}
	p.serviceFactories[name] = factory
	return nil
}

// Run launch all serices and wait to terminate
func (p *ServiceManager) RunAndWait() error {
	// Create all services
	for name, _ := range p.serviceFactories {
		quit := make(chan os.Signal)
		if service, err := p.serviceFactories[name].New(p.Config, quit); err != nil {
			glog.Infof("Create service '%s'failed", name)
			return err
		} else {
			glog.Infof("Create service '%s' successfully", name)
			p.Services[name] = service
			p.quits[name] = quit
		}
	}

	// Run all service
	glog.Infof("There are %d service in '%s'", len(p.Services), p.name)
	for _, service := range p.Services {
		glog.Infof("Starting service:'%s'...", service.Name())
		if err := service.Start(); err != nil {
			return err
		}
	}
	// Wait all service to terminate in main context
	for name, quit := range p.quits {
		<-quit
		glog.Info("Servide(%s) is terminated", name)
	}
	return nil
}

// Run launch all serices
func (p *ServiceManager) Run() error {
	// Run all service
	glog.Infof("There are %d service in '%s'", len(p.Services), p.name)
	for _, service := range p.Services {
		glog.Infof("Starting service:'%s'...", service.Name())
		if err := service.Start(); err != nil {
			return err
		}
	}
	return nil
}

// stop
func (p *ServiceManager) stop() error {
	for _, service := range p.Services {
		glog.Infof("Stoping service:'%s'...", service.Name())
		service.Stop()
	}
	return nil
}

/*
// StartService launch specified service
func (p *ServiceManager) StartService(name string) error {
	// Return error if service has already been started
	for id, service := range p.Services {
		if strings.IndexAny(id, name) >= 0 && service != nil {
			return fmt.Errorf("The service '%s' has already been started", name)
		}
	}
	quit := make(chan os.Signal)
	service, err := CreateService(name, p.Config, quit)
	if err != nil {
		glog.Errorf("%s", err)
	} else {
		glog.Infof("Create service '%s' success", name)
		p.Services[name] = service
		p.quits[name] = quit
	}
	return nil
}

// StopService stop specified service
func (p *ServiceManager) StopService(id string) error {
	for name, service := range p.Services {
		if name == id && service != nil {
			service.Stop()
			p.Services[name] = nil
			close(p.quits[name])
		}
	}
	return nil
}
*/
// GetServicesByName return service instance by name, or matched by part of name
func (p *ServiceManager) GetServicesByName(name string) []Service {
	services := []Service{}
	for k, service := range p.Services {
		if strings.IndexAny(k, name) >= 0 {
			services = append(services, service)
		}
	}
	return services
}

// GetService return service instance by name, or matched by part of name
func (p *ServiceManager) GetService(name string) Service {
	if _, found := p.Services[name]; found {
		return p.Services[name]
	}
	return nil
}
