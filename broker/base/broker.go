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
	"strings"

	"github.com/cloustone/sentel/core"
)

type HubNodeInfo struct {
	NodeName  string
	NodeIp    string
	CreatedAt string
}

// ServiceInfo
type ServiceInfo struct {
	ServiceName    string
	Listen         string
	Acceptors      uint64
	MaxClients     uint64
	CurrentClients uint64
	ShutdownCount  uint64
}
type Broker struct {
	core.ServiceManager
	nodeName string // Node name
}

const BrokerVersion = "0.1"

var (
	_broker *Broker
)

// GetBroker create service manager and all supported service
// The function should be called in service
func GetBroker() *Broker { return _broker }

// GetServicesByName return service instance by name, or matched by part of name
func (p *Broker) GetServicesByName(name string) []core.Service {
	services := []core.Service{}
	for k, service := range p.Services {
		if strings.IndexAny(k, name) >= 0 {
			services = append(services, service)
		}
	}
	return services
}

// GetAllProtocolServices() return all protocol services
func (p *Broker) GetAllProtocolServices() []ProtocolService {
	services := []ProtocolService{}
	for _, service := range p.Services {
		if p, ok := service.(ProtocolService); ok {
			services = append(services, p)
		}
	}
	return services
}

// GetProtocolServiceByname return protocol services by name
func (p *Broker) GetProtocolServices(name string) []ProtocolService {
	services := []ProtocolService{}
	for k, service := range p.Services {
		if strings.IndexAny(k, name) >= 0 {
			if p, ok := service.(ProtocolService); ok {
				services = append(services, p)
			}
		}
	}
	return services
}

// Node info
func (p *Broker) GetNodeInfo() *HubNodeInfo {
	return &HubNodeInfo{}
}

// Version
func (p *Broker) GetVersion() string {
	return BrokerVersion
}

// GetStats return server's stats
func (p *Broker) GetStats(serviceName string) map[string]uint64 {
	allstats := NewStats(false)
	services := p.GetProtocolServices(serviceName)

	for _, service := range services {
		stats := service.GetStats()
		allstats.AddStats(stats)
	}
	return allstats.Get()
}

// GetMetrics return server metrics
func (p *Broker) GetMetrics(serviceName string) map[string]uint64 {
	allmetrics := NewMetrics(false)
	services := p.GetProtocolServices(serviceName)

	for _, service := range services {
		metrics := service.GetMetrics()
		allmetrics.AddMetrics(metrics)
	}
	return allmetrics.Get()
}

// GetClients return clients list withspecified service
func (p *Broker) GetClients(serviceName string) []*ClientInfo {
	clients := []*ClientInfo{}
	services := p.GetProtocolServices(serviceName)

	for _, service := range services {
		list := service.GetClients()
		clients = append(clients, list...)
	}
	return clients
}

// GetClient return client info with specified client id
func (p *Broker) GetClient(serviceName string, id string) *ClientInfo {
	services := p.GetProtocolServices(serviceName)

	for _, service := range services {
		if client := service.GetClient(id); client != nil {
			return client
		}
	}
	return nil
}

// Kickoff Client killoff a client from specified service
func (p *Broker) KickoffClient(serviceName string, id string) error {
	services := p.GetProtocolServices(serviceName)

	for _, service := range services {
		if err := service.KickoffClient(id); err == nil {
			return nil
		}
	}
	return fmt.Errorf("Failed to kick off user '%s' from service '%s'", id, serviceName)
}

// GetSessions return all sessions information for specified service
func (p *Broker) GetSessions(serviceName string, conditions map[string]bool) []*SessionInfo {
	sessions := []*SessionInfo{}
	services := p.GetProtocolServices(serviceName)

	for _, service := range services {
		list := service.GetSessions(conditions)
		sessions = append(sessions, list...)
	}
	return sessions

}

// GetSession return specified session information with session id
func (p *Broker) GetSession(serviceName string, id string) *SessionInfo {
	services := p.GetProtocolServices(serviceName)

	for _, service := range services {
		if info := service.GetSession(id); info != nil {
			return info
		}
	}
	return nil
}

// GetRoutes return route table information for specified service
func (p *Broker) GetRoutes(serviceName string) []*RouteInfo {
	routes := []*RouteInfo{}
	services := p.GetProtocolServices(serviceName)

	for _, service := range services {
		list := service.GetRoutes()
		routes = append(routes, list...)
	}
	return routes
}

// GetRoute return route information for specified topic
func (p *Broker) GetRoute(serviceName string, topic string) *RouteInfo {
	services := p.GetProtocolServices(serviceName)

	for _, service := range services {
		route := service.GetRoute(topic)
		if route != nil {
			return route
		}
	}
	return nil
}

// GetTopics return topic informaiton for specified service
func (p *Broker) GetTopics(serviceName string) []*TopicInfo {
	topics := []*TopicInfo{}
	services := p.GetProtocolServices(serviceName)

	for _, service := range services {
		list := service.GetTopics()
		topics = append(topics, list...)
	}
	return topics
}

// GetTopic return topic information for specified topic
func (p *Broker) GetTopic(serviceName string, topic string) *TopicInfo {
	services := p.GetProtocolServices(serviceName)

	for _, service := range services {
		info := service.GetTopic(topic)
		if info != nil {
			return info
		}
	}
	return nil
}

// GetSubscriptions return subscription informaiton for specified service
func (p *Broker) GetSubscriptions(serviceName string) []*SubscriptionInfo {
	subs := []*SubscriptionInfo{}
	services := p.GetProtocolServices(serviceName)

	for _, service := range services {
		list := service.GetSubscriptions()
		subs = append(subs, list...)
	}
	return subs
}

// GetSubscription return subscription information for specified topic
func (p *Broker) GetSubscription(serviceName string, sub string) *SubscriptionInfo {
	services := p.GetProtocolServices(serviceName)

	for _, service := range services {
		info := service.GetSubscription(sub)
		if info != nil {
			return info
		}
	}
	return nil
}

// GetAllServiceInfo return all service information
func (p *Broker) GetAllServiceInfo() []*ServiceInfo {
	services := []*ServiceInfo{}

	//	for _, service := range p.Services {
	//		services = append(services, service.Info())
	//	}
	return services
}
