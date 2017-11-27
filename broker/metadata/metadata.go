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

package metadata

import "time"

type Device struct {
	DeviceId    string    `json:"deviceId" bson:"deviceId"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	LastUpdated time.Time `json:"lastUpdated" bson:"lastUpdated"`
	Data        []byte    `json:"data" bson:"data"`
}

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

// ServiceInfo
type ServiceInfo struct {
	ServiceName    string
	Listen         string
	Acceptors      uint64
	MaxClients     uint64
	CurrentClients uint64
	ShutdownCount  uint64
}

// SessionInfo
type SessionInfo struct {
	ClientId           string
	CleanSession       bool
	MessageMaxInflight uint64
	MessageInflight    uint64
	MessageInQueue     uint64
	MessageDropped     uint64
	AwaitingRel        uint64
	AwaitingComp       uint64
	AwaitingAck        uint64
	CreatedAt          string
}
