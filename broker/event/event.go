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
	SessionCreate    = 0x0001
	SessionDestroy   = 0x0002
	TopicPublish     = 0x0008
	TopicSubscribe   = 0x0010
	TopicUnsubscribe = 0x0100
	QutoChange       = 0x0200
	SessionResume    = 0x0400
	AuthChange       = 0x0800
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
	Topic     string `json:"topic"` // Topic
	Payload   []byte `json:"data"`  // Topic data
	Qos       uint8  `json:"qos"`
	Direction uint8  `json:"direction"`
	Retain    bool   `json:"retain"`
}

type QutoChangeType struct {
	QutoId string `json:"qutoId"`
}
