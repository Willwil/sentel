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
	Start() error
	Stop()
}

type ServiceFactory interface {
	New(c config.Config, quit chan os.Signal) (Service, error)
}

type ServiceManager struct {
	sync.Once
	Config           config.Config // Global config
	serviceFactories []ServiceFactory
	services         []Service                 // All service created by config.Protocols
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
func NewServiceManager(name string, c config.Config) (*ServiceManager, error) {
	if serviceManager != nil {
		return serviceManager, errors.New("NewServiceManager had been called many times")
	}
	mgr := &ServiceManager{
		Config:           c,
		serviceFactories: []ServiceFactory{},
		services:         []Service{},
		quits:            make(map[string]chan os.Signal),
		name:             name,
	}
	serviceManager = mgr
	return serviceManager, nil
}

func (p *ServiceManager) AddService(factory ServiceFactory) *ServiceManager {
	p.serviceFactories = append(p.serviceFactories, factory)
	return p
}

// Run launch all serices and wait to terminate
func (p *ServiceManager) RunAndWait() error {
	// Create all services
	for _, factory := range p.serviceFactories {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGUSR1)
		if service, err := factory.New(p.Config, quit); err != nil {
			return err
		} else {
			glog.Infof("Create service '%s' successfully", service.Name())
			p.services = append(p.services, service)
			p.quits[service.Name()] = quit
		}
	}

	// Run all service
	glog.Infof("There are %d services in '%s'", len(p.services), p.name)
	for _, service := range p.services {
		glog.Infof("Starting service:'%s'...", service.Name())
		if err := service.Start(); err != nil {
			return err
		}
	}
	// Wait all service to terminate in main context
	for name, quit := range p.quits {
		<-quit
		glog.Info("Service '%s' is terminated", name)
	}
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
		glog.Infof("Stoping service '%s'...", service.Name())
		service.Stop()
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
