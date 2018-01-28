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

package db

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

type ManagerDB struct {
	config  config.Config
	dbconn  *mgo.Session
	session *mgo.Database
}

const (
	dbnameIotmanager = "iotmanager"
)

func NewManagerDB(c config.Config) (*ManagerDB, error) {
	// try connect with mongo db
	addr := c.MustString("mongo")
	dbc, err := mgo.DialWithTimeout(addr, 1*time.Second)
	if err != nil {
		return nil, fmt.Errorf("connect with mongo '%s'failed: '%s'", addr, err.Error())
	}
	return &ManagerDB{
		config:  c,
		dbconn:  dbc,
		session: dbc.DB("iotmanager"),
	}, nil
}

func (p *ManagerDB) GetAllTenants() []Tenant {
	tenants := []Tenant{}
	c := p.session.C(collectionAdmin)
	c.Find(bson.M{}).All(&tenants)
	return tenants
}

func (p *ManagerDB) CreateTenant(t *Tenant) error {
	c := p.session.C(collectionAdmin)
	return c.Insert(t)
}

func (p *ManagerDB) RemoveTenant(tid string) error {
	c := p.session.C(collectionAdmin)
	return c.Remove(bson.M{"tenantId": tid})
}

func (p *ManagerDB) CreateProduct(tid string, pid string) error {
	c := p.session.C(collectionAdmin)
	pp := Product{ProductId: pid, CreatedAt: time.Now()}
	return c.Update(bson.M{"tenantId": tid},
		bson.M{"$addToSet": bson.M{"products": pp}})
}

func (p *ManagerDB) RemoveProduct(tid string, pid string) error {
	c := p.session.C(collectionAdmin)
	return c.Update(bson.M{"tenantId": tid},
		bson.M{"$pull": bson.M{"products": bson.M{"productId": pid}}})
}

func (p *ManagerDB) Close() {
	p.dbconn.Close()
}
