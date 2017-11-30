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

package event

import "github.com/cloustone/sentel/core"

const (
	SessionCreate    = 1
	SessionDestroy   = 2
	TopicPublish     = 3
	TopicSubscribe   = 4
	TopicUnsubscribe = 5
	QutoChange       = 6
	SessionResume    = 7
	AuthChange       = 8
)

type Event struct {
	core.TopicBase
	BrokerId string      `json:"brokerId"` // Broker identifier where event come from
	Type     uint32      `json:"type"`     // Event type
	ClientId string      `json:"clientId"` // Client identifier where event come from
	Detail   interface{} `json:"detail"`   // Event detail
}

type SessionCreateDetail struct {
	Persistent bool `json:"persistent"` // Whether the session is persistent
	Retain     bool `json:"retain"`
}

type SessionDestroyDetail struct {
	Persistent bool `json:"persistent"` // Whether the session is persistent
	Retain     bool `json:"retain"`
}

type TopicSubscribeDetail struct {
	Persistent bool   `json:"persistent"` // Whether the session is persistent
	Topic      string `json:"topic"`      // Topic
	Qos        uint8  `json:"qos"`
	Data       []byte `json:"data"` // Topic data
	Retain     bool   `json:"retain"`
}
type TopicUnsubscribeDetail struct {
	Persistent bool   `json:"persistent"` // Whether the session is persistent
	Topic      string `json:"topic"`      // Topic
	Data       []byte `json:"data"`       // Topic data
}

type TopicPublishDetail struct {
	Id        uint16 `json:"id"`
	Topic     string `json:"topic"` // Topic
	Payload   []byte `json:"data"`  // Topic data
	Qos       uint8  `json:"qos"`
	Direction uint8  `json:"direction"`
	Retain    bool   `json:"retain"`
}

type QutoChangeDetail struct {
	QutoId string `json:"qutoId"`
}
