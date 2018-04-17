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

package mgrdb

import (
	"errors"
	"time"

	"github.com/cloustone/sentel/pkg/config"
)

const (
	DBNAME = "iotmanager"
)

// Node
type Node struct {
	NodeId     string    `json:"nodeId" bson:"NodeId"`
	NodeIp     string    `json:"nodeIp" bson:"NodeIp"`
	Version    string    `json:"version" bson:"Version"`
	CreatedAt  time.Time `json:"createdAt" bson:"CreatedAt"`
	NodeStatus string    `json:"nodeStatus" bson:"NodeStatus"`
	UpdatedAt  time.Time `json:"updatedAt" bson:"UpdatedAt"`
	Action     string    `json:"action"`
}

// Subscription
type Subscription struct {
	ClientId        string    `json:"clientID" bson:"ClientId"`
	SubscribedTopic string    `json:"topic" bson:"Topic"`
	Qos             int       `json:"qos" bson:"Qos"`
	CreatedAt       time.Time `json:"createdAt" bson:"CreatedAt"`
	Action          string    `json:"action" bson:"Action"`
}

// Client
type Client struct {
	ClientId        string `json:"clientID" bson:"ClientId"`
	UserName        string `json:"userName" bson:"UserName"`
	IpAddress       string `json:"ipAddress" bson:"IpAddress"`
	Port            uint16 `json:"port" bson:"Port"`
	CleanSession    bool   `json:"cleanSession" bson:"CleanSession"`
	ProtocolVersion string `json:"protocolVersion" bson:"ProtocolVersion"`
	Keepalive       uint16 `json:"keepalive" bson:"Keepalive"`
	ConnectedAt     string `json:"connectedAt" bson:"ConnectedAt"`
}

// Session
type Session struct {
	Action             string `json:"action" bson:"Action"`
	ClientId           string `json:"clientID" bson:"ClientId"`
	CleanSession       bool   `json:"cleanSession" bson:"CleanSession"`
	MessageMaxInflight uint64 `json:"messageMaxInflight" bson:"MessageMaxInflight"`
	MessageInflight    uint64 `json:"messageInflight" bson:"MessageInflight"`
	MessageInQueue     uint64 `json:"messageInQueue" bson:"MessageInQueue"`
	MessageDropped     uint64 `json:"messageDropped" bson:"MessageDropped"`
	AwaitingRel        uint64 `json:"awaitingRel" bson:"AwaitingRel"`
	AwaitingComp       uint64 `json:"awaitingComp" bson:"AwaitingComp"`
	AwaitingAck        uint64 `json:"awaitingAck" bson:"AwaitingAck"`
	CreatedAt          string `json:"createdAt" bson:"CreatedAt"`
}

// Tenant & Product
type Tenant struct {
	TenantId         string              `bson:"tenantId"`
	Products         map[string]*Product `bson:"products"`
	ServiceName      string              `bson:"serviceName"`
	ServiceId        string              `bson:"serviceId"`
	ServiceState     string              `bson:"servcieState"`
	InstanceReplicas int32               `bson:"instanceReplicas"`
	NetworkId        string              `bson:"networkId"`
	CreatedAt        time.Time           `bson:"createdAt"`
}

type Product struct {
	ProductId string    `bson:"productId"`
	CreatedAt time.Time `bson:"createdAt"`
}

// Metric
type Metric struct {
	Action     string            `json:"action" bson:"Action"`
	NodeName   string            `json:"nodeName" bson:"NodeName"`
	Service    string            `json:"service" bson:"Service"`
	Values     map[string]uint64 `json:"values" bson:"Values"`
	UpdateTime time.Time         `json:"updateTime" bson:"UpdateTime"`
}

// Stat
type Stats struct {
	NodeName   string            `json:"nodeName" bson:"NodeName"`
	Service    string            `json:"service" bson:"Service"`
	Action     string            `json:"action" bson:"Action"`
	UpdateTime time.Time         `json:"updateTime" bson:"UpdateTime"`
	Values     map[string]uint64 `json:"values" bson:"Values"`
}

// Publish
type Publish struct {
	ClientId        string `json:"clientID" bson:"ClientId"`
	SubscribedTopic string `json:"topic" bson:"SubscribedTopic"`
	ProductId       string `json:"product" bson:"ProductId"`
}

type ManagerDB interface {
	// Node
	// GetAllNodes retrieve all broker nodes info in cluster
	GetAllNodes() []Node
	// GetNode retrieve specified node info
	GetNode(nodeId string) (Node, error)
	// AddNode add new node
	AddNode(n Node) error
	// UpdateNode update existed node
	UpdateNode(n Node) error
	// RemoveNode remove existed node
	RemoveNode(nodeId string) error
	// GetNodeClients retrieve all clients info on specified broker node
	GetNodeClients(nodeId string) []Client
	// GetNodeClients retrieve clients info on specified broker node in duration
	GetNodesClientWithTimeScope(nodeId string, from time.Time, to time.Time) []Client
	// GetClientWithNode return specified client info on specified broker node
	GetClientWithNode(nodeId string, clientID string) (Client, error)

	// Session
	// GetSession return client's session info
	GetSession(clientID string) (Session, error)
	// UpdateSession update session info
	UpdateSession(s Session) error
	// GetNodeSessions retrieve all sessions on specified broker node
	GetNodeSessions(nodeId string) []Session
	// GetSessionWithNode retrieve specifed client session on specifed broker node
	GetSessionWithNode(nodeId string, clientID string) (Session, error)

	// Tenants
	// GetAllTenants retrieve all tenants in cluster
	GetAllTenants() []Tenant
	// AddTenant add tenant into db
	AddTenant(t *Tenant) error
	// RemoveTenant remove tenant from db
	RemoveTenant(tenantId string) error

	// Product
	AddProduct(tenantId string, productId string) error
	RemoveProduct(tenantId string, productId string) error

	// Client
	GetClient(clientID string) (Client, error)
	AddClient(Client) error
	RemoveClient(Client) error
	UpdateClient(Client) error

	// Metrics
	GetMetrics() []Metric
	GetNodeMetric(nodeId string) (Metric, error)
	AddMetricHistory(Metric) error
	UpdateMetric(Metric) error

	// Stats
	AddStatsHistory(Stats) error
	UpdateStats(Stats) error
	GetStats() []Stats

	// Subscription
	AddSubscription(sub Subscription) error
	UpdateSubscription(sub Subscription) error
	RemoveSubscription(sub Subscription) error

	// Publish
	UpdatePublish(Publish) error

	Close()
}

func New(c config.Config) (ManagerDB, error) {
	if _, err := c.String("mongo"); err == nil {
		return newMgrdbMongo(c)
	}
	return nil, errors.New("no valid database")
}
