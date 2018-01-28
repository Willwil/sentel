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

// Metric
type Metric struct {
	Action     string            `json:"action" bson:"Action"`
	NodeName   string            `json:"nodeName" bson:"NodeName"`
	Service    string            `json:"service" bson:"Service"`
	Values     map[string]uint64 `json:"values" bson:"Values"`
	UpdateTime time.Time         `json:"updateTime" bson:"UpdateTime"`
}

func (p *ManagerDB) GetMetrics() []Metric {
	c := p.session.C(collectionMetrics)
	metrics := []Metric{}
	c.Find(nil).Iter().All(&metrics)
	return metrics
}

func (p *ManagerDB) GetNodeMetric(nodeId string) (Metric, error) {
	c := p.session.C(collectionMetrics)
	metric := Metric{}
	err := c.Find(bson.M{"NodeId": nodeId}).One(&metric)
	return metric, err
}
