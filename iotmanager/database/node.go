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
//  under the License.

package db

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Node
type Node struct {
	TopicName  string
	NodeId     string    `json:"nodeId" bson:"NodeId"`
	NodeIp     string    `json:"nodeIp" bson:"NodeIp"`
	Version    string    `json:"version" bson:"Version"`
	CreatedAt  time.Time `json:"createdAt" bson:"CreatedAt"`
	NodeStatus string    `json:"nodeStatus" bson:"NodeStatus"`
	UpdatedAt  time.Time `json:"updatedAt" bson:"UpdatedAt"`
	Action     string    `json:"action"`
}

func (p *ManagerDB) GetAllNodes() []Node {
	c := p.session.C(collectionNodes)
	nodes := []Node{}
	c.Find(nil).All(&nodes)
	return nodes
}

func (p *ManagerDB) GetNode(nodeId string) (Node, error) {
	c := p.session.C(collectionNodes)
	node := Node{}
	err := c.Find(bson.M{"NodeId": nodeId}).One(&node)
	return node, err
}

func (p *ManagerDB) GetNodeClients(nodeId string) []Client {
	c := p.session.C(collectionClients)
	clients := []Client{}
	c.Find(bson.M{"NodeId": nodeId}).Iter().All(&clients)
	return clients
}

func (p *ManagerDB) GetNodesClientWithTimeScope(nodeId string, form time.Time, to time.Time) []Client {
	clients := []Client{}
	// TODO
	return clients
}

func (p *ManagerDB) GetClientWithNode(nodeId string, clientId string) (Client, error) {
	client := Client{}
	// TODO
	return client, nil
}
