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

package cluster

import "time"

type tenant struct {
	tid       string    `json:"tenantId"`
	createdAt time.Time `json:"createdAt"`
	products  map[string]*product
}

type product struct {
	pid       string    `json:"productId"`
	tid       string    `json:"tenantId"`
	createdAt time.Time `json:"createdAt"`
	brokers   map[string]*broker
}

type broker struct {
	bid         string       // broker identifier
	tid         string       // tenant identifier
	ip          string       // broker ip address
	port        string       // broker port
	status      brokerStatus // broker status
	createdAt   time.Time    // created time for broker
	lastUpdated time.Time    // last updated time for broker
	context     interface{}  // the attached context
}

const (
	brokerStatusInvalid = "invalid"
	brokerStatusStarted = "started"
	brokerStatusStoped  = "stoped"
)

type brokerStatus string
