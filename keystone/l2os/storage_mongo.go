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
package l2

import (
	"fmt"
	"time"

	"github.com/cloustone/sentel/pkg/config"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoStorage struct {
	config  config.Config
	session *mgo.Session
	db      *mgo.Database
}

func newMongoStorage(c config.Config) (*mongoStorage, error) {
	hosts := c.MustString("keystone", "mongo")
	timeout := c.MustInt("keystone", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}

	db := session.DB("keystone")
	return &mongoStorage{config: c, session: session, db: db}, nil
}

func (p *mongoStorage) getAccount(cname string) (*Account, error) {
	account := Account{}
	c := p.db.C("accounts")
	if err := c.Find(bson.M{"name": cname}).One(&account); err != nil {
		return nil, fmt.Errorf("account '%s' no exist", account.Name)
	}
	return &account, nil
}

func (p *mongoStorage) createObject(obj *Object) error {
	c := p.db.C("objects")
	b := Object{}
	if err := c.Find(bson.M{"objectId": obj.ObjectId}).One(&b); err == nil {
		return fmt.Errorf("object '%s' already exist", obj.ObjectId)
	}
	return c.Insert(obj)
}

func (p *mongoStorage) deleteObject(objid string) error {
	c := p.db.C("objects")
	return c.Remove(bson.M{"objectId": objid})
}

func (p *mongoStorage) getObject(objid string) (*Object, error) {
	c := p.db.C("objects")
	obj := Object{}
	if err := c.Find(bson.M{"objectId": objid}).One(&obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

func (p *mongoStorage) updateObject(obj *Object) error {
	c := p.db.C("objects")
	return c.Update(bson.M{"objectId": obj.ObjectId}, obj)
}

func (p *mongoStorage) deinitialize() {
	if p.session != nil {
		p.session.Close()
	}
}

func (p *mongoStorage) createAccount(name string) error {
	c := p.db.C("accounts")
	if err := c.Find(bson.M{"name": name}); err == nil {
		return fmt.Errorf("account '%s' already exist", name)
	}
	return c.Insert(bson.M{"name": name, "createdAt": time.Now()})
}
func (p *mongoStorage) destroyAccount(name string) error {
	c := p.db.C("accounts")
	err := c.Remove(bson.M{"name": name})
	if err == nil {
		c := p.db.C(name)
		return c.DropCollection()
	}
	return err
}
