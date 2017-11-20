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

// ClientInfo
type ClientInfo struct {
	UserName     string
	CleanSession bool
	PeerName     string
	ConnectTime  string
}

// RouteInfo
type RouteInfo struct {
	Topic string
	Route []string
}

// TopicInfo
type TopicInfo struct {
	Topic     string
	Attribute string
}

// SubscriptionInfo
type SubscriptionInfo struct {
	ClientId  string
	Topic     string
	Attribute string
}

type ProtocolService interface {
	// Service Infomation
	GetServiceInfo() *ServiceInfo

	// Client
	GetClients() []*ClientInfo
	GetClient(id string) *ClientInfo
	KickoffClient(id string) error

	// Session
	GetSessions(conditions map[string]bool) []*SessionInfo
	GetSession(id string) *SessionInfo

	// Route info
	GetRoutes() []*RouteInfo
	GetRoute(id string) *RouteInfo

	// Topic info
	GetTopics() []*TopicInfo
	GetTopic(id string) *TopicInfo

	// SubscriptionInfo
	GetSubscriptions() []*SubscriptionInfo
	GetSubscription(id string) *SubscriptionInfo
}
