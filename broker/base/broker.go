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

package base

import (
	"fmt"

	"github.com/cloustone/sentel/core"
	uuid "github.com/satori/go.uuid"
)

type Broker struct {
	core.ServiceManager
	brokerId string
}

const BrokerVersion = "0.1"

var (
	_broker *Broker
)

// newBroker create global broker
func NewBroker(c core.Config) (*Broker, error) {
	if _broker != nil {
		panic("Global broker had already been created")
	}
	serviceMgr, err := core.NewServiceManager("broker", c)
	if err != nil {
		return nil, err
	}
	_broker = &Broker{
		ServiceManager: *serviceMgr,
		brokerId:       uuid.NewV4().String(),
	}
	return _broker, nil
}

// GetBroker create service manager and all supported service
// The function should be called in service
func GetBroker() *Broker { return _broker }

// Version
func GetBrokerVersion() string {
	return BrokerVersion
}

// GetBrokerId return broker's identifier
func GetBrokerId() string {
	return _broker.brokerId
}

// GetService return specified service instance
func GetService(name string) core.Service {
	broker := GetBroker()
	return broker.GetServiceByName(name)
}

// GetConfig return broker's configuration
func (p *Broker) GetConfig() core.Config {
	return p.Config
}

// GetServicesByName return service instance by name, or matched by part of name
func (p *Broker) GetServiceByName(name string) core.Service {
	if _, ok := p.Services[name]; !ok {
		panic(fmt.Sprintf("Failed to find service '%s' in broker", name))
	}
	return p.Services[name]
}
