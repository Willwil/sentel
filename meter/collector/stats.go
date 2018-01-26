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

package collector

import (
	"context"
	"time"

	"github.com/cloustone/sentel/pkg/message"

	"gopkg.in/mgo.v2/bson"
)

// Stat
type Stats struct {
	TopicName  string
	NodeName   string            `json:"nodeName"`
	Service    string            `json:"service"`
	Action     string            `json:"action"`
	UpdateTime time.Time         `json:"updateTime"`
	Values     map[string]uint64 `json:"values"`
}

func (p *Stats) Topic() string        { return TopicNameStats }
func (p *Stats) SetTopic(name string) {}
func (p *Stats) Serialize(opt message.SerializeOption) ([]byte, error) {
	return message.Serialize(p, opt)
}
func (p *Stats) Deserialize(buf []byte, opt message.SerializeOption) error { return nil }

func (p *Stats) handleTopic(service *collectorService, ctx context.Context) error {
	db, err := service.getDatabase()
	if err != nil {
		return err
	}
	defer db.Session.Close()

	switch p.Action {
	case ObjectActionUpdate:
		// update newest stats for the node
		c := db.C("stats")
		node := Node{}
		if err := c.Find(bson.M{"NodeName": p.NodeName}).One(&node); err != nil { // not found
			c.Insert(&Stats{
				NodeName:   p.NodeName,
				Service:    p.Service,
				Values:     p.Values,
				UpdateTime: time.Now(),
			})
		} else {
			c.Update(&node,
				&Stats{
					NodeName:   p.NodeName,
					Service:    p.Service,
					Values:     p.Values,
					UpdateTime: time.Now(),
				})
		}
		// save history data
		c = db.C("stats_history")
		c.Insert(&Stats{
			NodeName:   p.NodeName,
			Service:    p.Service,
			Values:     p.Values,
			UpdateTime: time.Now(),
		})
	case ObjectActionDelete:
	case ObjectActionRegister:
	default:
	}
	return nil
}
