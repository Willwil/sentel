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

	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

type mongoQueuePlugin struct {
	id             string
	config         core.Config
	hosts          string
	connectTimeout time.Duration
}

func newMongoQueuePlugin(id string, c core.Config) (queuePlugin, error) {
	// check mongo db configuration
	hosts, _ := core.GetServiceEndpoint(c, "broker", "mongo")
	timeout := c.MustInt("broker", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	return &mongoQueuePlugin{
		id:             id,
		config:         c,
		hosts:          hosts,
		connectTimeout: time.Duration(timeout) * time.Second,
	}, nil
}

func (p *mongoQueuePlugin) getData() (*queueData, error) {
	session, err := mgo.DialWithTimeout(p.hosts, p.connectTimeout)
	if err != nil {
		glog.Fatal("queue: failed to connect with backend data queue")
		return nil, err
	}
	defer session.Close()
	c := session.DB("persistent-session").C(p.id)
	data := queueData{}
	err = c.Find(nil).One(&data)
	return &data, nil

}
func (p *mongoQueuePlugin) pushData(data *queueData) {
	session, err := mgo.DialWithTimeout(p.hosts, p.connectTimeout)
	if err != nil {
		glog.Fatal("queue: failed to connect with backend data queue")
		return
	}
	defer session.Close()
	c := session.DB("persistent-session").C(p.id)
	c.Insert(data)
}
