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

package core

import (
	"errors"
	"fmt"
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

var (
	_serviceFactories = make(map[string]ServiceFactory)
)

// RegisterService register service with name and protocol specified
func RegisterService(name string, factory ServiceFactory) {
	if _, ok := _serviceFactories[name]; ok {
		glog.Errorf("Service '%s' is not registered", name)
	}
	_serviceFactories[name] = factory
}

// RegisterServiceWithConfig register service with name and configuration
func RegisterServiceWithConfig(name string, factory ServiceFactory, configs map[string]string) {
	if _, ok := _serviceFactories[name]; ok {
		glog.Errorf("Service '%s' is not registered", name)
	}
	RegisterConfig(name, configs)
	_serviceFactories[name] = factory
}

// CreateService create service instance according to service name
func CreateService(name string, c Config, ch chan os.Signal) (Service, error) {
	glog.Infof("Creating service '%s'...", name)

	if factory, ok := _serviceFactories[name]; ok && factory != nil {
		return _serviceFactories[name].New(c, ch)
	}
	return nil, fmt.Errorf("Invalid service '%s'", name)
}

// RunWithConfigFile create and start server
func RunWithConfigFile(serverName string, fileName string) error {
	glog.Infof("Starting '%s' server...", serverName)

	// Check all registered service
	if err := CheckAllRegisteredServices(); err != nil {
		return err
	}
	// Get configuration
	config, err := NewWithConfigFile(fileName)
	if err != nil {
		return err
	}
	// Create service manager according to the configuration
	mgr, err := NewServiceManager(serverName, config)
	if err != nil {
		return err
	}
	return mgr.Run()
}

// CheckAllRegisteredServices check all registered service simplily
func CheckAllRegisteredServices() error {
	if len(_serviceFactories) == 0 {
		return errors.New("No service registered")
	}
	for name, _ := range _serviceFactories {
		glog.Infof("Service '%s' is registered", name)
	}
	return nil
}

type ServiceManager struct {
	sync.Once
	Config   Config                    // Global config
	Services map[string]Service        // All service created by config.Protocols
	quits    map[string]chan os.Signal // Notification channel for each service
	name     string                    // Name of service manager
}

const serviceManagerVersion = "0.1"

var (
	_serviceManager *ServiceManager
)

// GetServiceManager create service manager and all supported service
// The function should be called in service
func GetServiceManager() *ServiceManager { return _serviceManager }

// NewServiceManager create ServiceManager only in main context
func NewServiceManager(name string, c Config) (*ServiceManager, error) {
	if _serviceManager != nil {
		return _serviceManager, errors.New("NewServiceManager had been called many times")
	}
	mgr := &ServiceManager{
		Config:   c,
		quits:    make(map[string]chan os.Signal),
		Services: make(map[string]Service),
		name:     name,
	}
	// Get supported configs
	items := c.MustString(name, "services")
	services := strings.Split(items, ",")
	// Create service for each protocol
	for _, name := range services {
		// Format service name
		name = strings.Trim(name, " ")
		quit := make(chan os.Signal)
		service, err := CreateService(name, c, quit)
		if err != nil {
			// If one of service is not started, must return
			return nil, err
		} else {
			glog.Infof("Create service '%s' successfully", name)
			mgr.Services[name] = service
			mgr.quits[name] = quit
		}
	}
	_serviceManager = mgr
	return _serviceManager, nil
}

// Run launch all serices and wait to terminate
func (p *ServiceManager) Run() error {
	// Run all service
	glog.Infof("There are %d service in '%s'", len(p.Services), p.name)
	for _, service := range p.Services {
		glog.Infof("Starting service:'%s'...", service.Name())
		go service.Start()
	}
	// Wait all service to terminate in main context
	for name, quit := range p.quits {
		<-quit
		glog.Info("Servide(%s) is terminated", name)
	}
	return nil
}

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
