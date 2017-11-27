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

package broker

import "github.com/cloustone/sentel/core"

// Publish publish event to event service
func Notify(e *Event) {
	e.BrokerId = GetId()
	broker := GetBroker()
	broker.eventmgr.notify(e)
}

// Subscribe subcribe event from event service
func Subscribe(event uint32, handler EventHandler, ctx interface{}) {
	broker := GetBroker()
	broker.eventmgr.subscribe(event, handler, ctx)
}

// GetBroker create service manager and all supported service
// The function should be called in service
func GetBroker() *Broker { return broker }

// Version
func GetVersion() string {
	return BrokerVersion
}

// GetBrokerId return broker's identifier
func GetId() string {
	return broker.brokerId
}

// GetService return specified service instance
func GetService(name string) core.Service {
	return broker.getServiceByName(name)
}

// GetConfig return broker's configuration
func (p *Broker) GetConfig() core.Config {
	return p.config
}
