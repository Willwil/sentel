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
	"fmt"
	"time"

	"github.com/cloustone/sentel/pkg/config"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	collectionNodes    = "nodes"
	collectionClients  = "clients"
	collectionMetrics  = "metrics"
	collectionAdmin    = "admin"
	collectionSessions = "sessions"
)

type mgrdbMongo struct {
	config  config.Config
	dbconn  *mgo.Session
	session *mgo.Database
}

func newMgrdbMongo(c config.Config) (ManagerDB, error) {
	// try connect with mongo db
	addr := c.MustString("mongo")
	dbc, err := mgo.DialWithTimeout(addr, 1*time.Second)
	if err != nil {
		return nil, fmt.Errorf("connect with mongo '%s'failed: '%s'", addr, err.Error())
	}
	return &mgrdbMongo{
		config:  c,
		dbconn:  dbc,
		session: dbc.DB(DBNAME),
	}, nil
}

func (p *mgrdbMongo) Close() {
	p.dbconn.Close()
}

// Node
func (p *mgrdbMongo) GetAllNodes() []Node {
	c := p.session.C(collectionNodes)
	nodes := []Node{}
	c.Find(nil).All(&nodes)
	return nodes
}

func (p *mgrdbMongo) GetNode(nodeId string) (Node, error) {
	c := p.session.C(collectionNodes)
	node := Node{}
	err := c.Find(bson.M{"NodeId": nodeId}).One(&node)
	return node, err
}

func (p *mgrdbMongo) GetNodeClients(nodeId string) []Client {
	c := p.session.C(collectionClients)
	clients := []Client{}
	c.Find(bson.M{"NodeId": nodeId}).Iter().All(&clients)
	return clients
}

func (p *mgrdbMongo) GetNodesClientWithTimeScope(nodeId string, form time.Time, to time.Time) []Client {
	clients := []Client{}
	// TODO
	return clients
}

func (p *mgrdbMongo) GetClientWithNode(nodeId string, clientId string) (Client, error) {
	client := Client{}
	// TODO
	return client, nil
}

// Session
func (p *mgrdbMongo) GetSession(clientId string) (Session, error) {
	c := p.session.C(collectionSessions)
	session := Session{}
	err := c.Find(bson.M{"ClientId": clientId}).One(&session)
	return session, err
}

func (p *mgrdbMongo) UpdateSession(s Session) error {
	c := p.session.C(collectionSessions)
	session := Session{}
	if err := c.Find(bson.M{"ClientId": s.ClientId}).One(&session); err != nil {
		return c.Insert(s)
	} else {
		return c.Update(session, s)
	}
}

func (p *mgrdbMongo) GetNodeSessions(nodeId string) []Session {
	c := p.session.C(collectionSessions)
	sessions := []Session{}
	c.Find(bson.M{"NodeId": nodeId}).Iter().All(&sessions)
	return sessions
}

func (p *mgrdbMongo) GetSessionWithNode(nodeId string, clientId string) (Session, error) {
	c := p.session.C(collectionSessions)
	session := Session{}
	err := c.Find(bson.M{"NodeId": nodeId, "ClientId": clientId}).One(&session)
	return session, err
}

func (p *mgrdbMongo) GetAllTenants() []Tenant {
	tenants := []Tenant{}
	c := p.session.C(collectionAdmin)
	c.Find(bson.M{}).All(&tenants)
	return tenants
}

func (p *mgrdbMongo) AddTenant(t *Tenant) error {
	c := p.session.C(collectionAdmin)
	return c.Insert(t)
}

func (p *mgrdbMongo) RemoveTenant(tid string) error {
	c := p.session.C(collectionAdmin)
	return c.Remove(bson.M{"tenantId": tid})
}

func (p *mgrdbMongo) AddProduct(tid string, pid string) error {
	c := p.session.C(collectionAdmin)
	pp := Product{ProductId: pid, CreatedAt: time.Now()}
	return c.Update(bson.M{"tenantId": tid},
		bson.M{"$addToSet": bson.M{"products": pp}})
}

func (p *mgrdbMongo) RemoveProduct(tid string, pid string) error {
	c := p.session.C(collectionAdmin)
	return c.Update(bson.M{"tenantId": tid},
		bson.M{"$pull": bson.M{"products": bson.M{"productId": pid}}})
}

func (p *mgrdbMongo) GetClient(clientId string) (Client, error) {
	c := p.session.C(collectionClients)
	client := Client{}
	err := c.Find(bson.M{"ClientId": clientId}).One(&client)
	return client, err
}

func (p *mgrdbMongo) UpdateClient(client Client) error {
	c := p.session.C(collectionClients)
	result := Client{}
	if err := c.Find(bson.M{"ClientId": client.ClientId}).One(&result); err == nil {
		return c.Update(result, p)
	} else {
		return c.Insert(p)
	}
}

func (p *mgrdbMongo) GetMetrics() []Metric {
	c := p.session.C(collectionMetrics)
	metrics := []Metric{}
	c.Find(nil).Iter().All(&metrics)
	return metrics
}

func (p *mgrdbMongo) GetNodeMetric(nodeId string) (Metric, error) {
	c := p.session.C(collectionMetrics)
	metric := Metric{}
	err := c.Find(bson.M{"NodeId": nodeId}).One(&metric)
	return metric, err
}
