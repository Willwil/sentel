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
func (this *Broker) GetServicesByName(name string) []core.Service {
	services := []core.Service{}
	for k, service := range this.Services {
		if strings.IndexAny(k, name) >= 0 {
			services = append(services, service)
		}
	}
	return services
}

// GetAllProtocolServices() return all protocol services
func (this *Broker) GetAllProtocolServices() []ProtocolService {
	services := []ProtocolService{}
	for _, service := range this.Services {
		if p, ok := service.(ProtocolService); ok {
			services = append(services, p)
		}
	}
	return services
}

// GetProtocolServiceByname return protocol services by name
func (this *Broker) GetProtocolServices(name string) []ProtocolService {
	services := []ProtocolService{}
	for k, service := range this.Services {
		if strings.IndexAny(k, name) >= 0 {
			if p, ok := service.(ProtocolService); ok {
				services = append(services, p)
			}
		}
	}
	return services
}

// Node info
func (this *Broker) GetNodeInfo() *HubNodeInfo {
	return &HubNodeInfo{}
}

// Version
func (this *Broker) GetVersion() string {
	return BrokerVersion
}

// GetStats return server's stats
func (this *Broker) GetStats(serviceName string) map[string]uint64 {
	allstats := NewStats(false)
	services := this.GetProtocolServices(serviceName)

	for _, service := range services {
		stats := service.GetStats()
		allstats.AddStats(stats)
	}
	return allstats.Get()
}

// GetMetrics return server metrics
func (this *Broker) GetMetrics(serviceName string) map[string]uint64 {
	allmetrics := NewMetrics(false)
	services := this.GetProtocolServices(serviceName)

	for _, service := range services {
		metrics := service.GetMetrics()
		allmetrics.AddMetrics(metrics)
	}
	return allmetrics.Get()
}

// GetClients return clients list withspecified service
func (this *Broker) GetClients(serviceName string) []*ClientInfo {
	clients := []*ClientInfo{}
	services := this.GetProtocolServices(serviceName)

	for _, service := range services {
		list := service.GetClients()
		clients = append(clients, list...)
	}
	return clients
}

// GeteClient return client info with specified client id
func (this *Broker) GetClient(serviceName string, id string) *ClientInfo {
	services := this.GetProtocolServices(serviceName)

	for _, service := range services {
		if client := service.GetClient(id); client != nil {
			return client
		}
	}
	return nil
}

// Kickoff Client killoff a client from specified service
func (this *Broker) KickoffClient(serviceName string, id string) error {
	services := this.GetProtocolServices(serviceName)

	for _, service := range services {
		if err := service.KickoffClient(id); err == nil {
			return nil
		}
	}
	return fmt.Errorf("Failed to kick off user '%s' from service '%s'", id, serviceName)
}

// GetSessions return all sessions information for specified service
func (this *Broker) GetSessions(serviceName string, conditions map[string]bool) []*SessionInfo {
	sessions := []*SessionInfo{}
	services := this.GetProtocolServices(serviceName)

	for _, service := range services {
		list := service.GetSessions(conditions)
		sessions = append(sessions, list...)
	}
	return sessions

}

// GetSession return specified session information with session id
func (this *Broker) GetSession(serviceName string, id string) *SessionInfo {
	services := this.GetProtocolServices(serviceName)

	for _, service := range services {
		if info := service.GetSession(id); info != nil {
			return info
		}
	}
	return nil
}

// GetRoutes return route table information for specified service
func (this *Broker) GetRoutes(serviceName string) []*RouteInfo {
	routes := []*RouteInfo{}
	services := this.GetProtocolServices(serviceName)

	for _, service := range services {
		list := service.GetRoutes()
		routes = append(routes, list...)
	}
	return routes
}

// GetRoute return route information for specified topic
func (this *Broker) GetRoute(serviceName string, topic string) *RouteInfo {
	services := this.GetProtocolServices(serviceName)

	for _, service := range services {
		route := service.GetRoute(topic)
		if route != nil {
			return route
		}
	}
	return nil
}

// GetTopics return topic informaiton for specified service
func (this *Broker) GetTopics(serviceName string) []*TopicInfo {
	topics := []*TopicInfo{}
	services := this.GetProtocolServices(serviceName)

	for _, service := range services {
		list := service.GetTopics()
		topics = append(topics, list...)
	}
	return topics
}

// GetTopic return topic information for specified topic
func (this *Broker) GetTopic(serviceName string, topic string) *TopicInfo {
	services := this.GetProtocolServices(serviceName)

	for _, service := range services {
		info := service.GetTopic(topic)
		if info != nil {
			return info
		}
	}
	return nil
}

// GetSubscriptions return subscription informaiton for specified service
func (this *Broker) GetSubscriptions(serviceName string) []*SubscriptionInfo {
	subs := []*SubscriptionInfo{}
	services := this.GetProtocolServices(serviceName)

	for _, service := range services {
		list := service.GetSubscriptions()
		subs = append(subs, list...)
	}
	return subs
}

// GetSubscription return subscription information for specified topic
func (this *Broker) GetSubscription(serviceName string, sub string) *SubscriptionInfo {
	services := this.GetProtocolServices(serviceName)

	for _, service := range services {
		info := service.GetSubscription(sub)
		if info != nil {
			return info
		}
	}
	return nil
}

// GetAllServiceInfo return all service information
func (this *Broker) GetAllServiceInfo() []*ServiceInfo {
	services := []*ServiceInfo{}

	//	for _, service := range this.Services {
	//		services = append(services, service.Info())
	//	}
	return services
}
