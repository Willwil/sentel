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

package queue

import (
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/pkg/config"
)

type mongoPlugin struct {
	clientId   string
	config     config.Config
	session    *mgo.Session
	collection *mgo.Collection
}

func newMongoPlugin(clientId string, c config.Config) (queuePlugin, error) {
	// check mongo db configuration
	hosts := c.MustString("broker", "mongo")
	timeout := c.MustInt("broker", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	return &mongoPlugin{
		clientId:   clientId,
		config:     c,
		session:    session,
		collection: session.DB("persistent-session").C(clientId),
	}, nil
}

// Initialize initialize the backend queue plugin
func (p *mongoPlugin) initialize() error {
	return nil
}

// GetMessageCount return total mesasge count in queue
func (p *mongoPlugin) length() int {
	if count, err := p.collection.Find(nil).Count(); err == nil {
		return count
	}
	return -1
}

// Front return message at queue's head
func (p *mongoPlugin) front() *base.Message {
	msg := base.Message{}
	if err := p.collection.Find(nil).One(&msg); err != nil {
		return nil
	}
	return &msg
}

// Pushback push message at tail of queue
func (p *mongoPlugin) pushback(msg *base.Message) {
	p.collection.Insert(msg)
}

// Pop popup head message from queue
func (p *mongoPlugin) pop() *base.Message {
	msg := base.Message{}
	if err := p.collection.Find(nil).One(&msg); err == nil {
		p.collection.Remove(&msg)
		return &msg
	}
	return nil
}

// Close closes the connection.
func (p *mongoPlugin) close() error {
	p.session.Close()
	return nil
}
