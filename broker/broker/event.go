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

package broker

import "github.com/cloustone/sentel/core"

const (
	SessionCreated    = 0x0001
	SessionDestroyed  = 0x0002
	TopicPublished    = 0x0008
	TopicSubscribed   = 0x0010
	TopicUnsubscribed = 0x0100
	QutoChanged       = 0x0200
	SessionResumed    = 0x0400
	AuthChanged       = 0x0800
)

type Event struct {
	core.TopicBase
	BrokerId   string `json:"brokerId"`   // Broker identifier where event come from
	Type       uint32 `json:"type"`       // Event type
	ClientId   string `json:"clientId"`   // Client identifier where event come from
	Topic      string `json:"topic"`      // Topic
	Data       []byte `json:"data"`       // Topic data
	Persistent bool   `json:"persistent"` // Whether the session is persistent
	QutoId     string `json:"qutoId"`     // Qutotation identifier
	Qos        uint8  `json:"qos"`
}

type SessionCreateEvent struct {
	Event
	Persistent bool `json:"persistent"` // Whether the session is persistent
}

type TopicSubscribeEvent struct {
	Event
	Topic string `json:"topic"` // Topic
	Data  []byte `json:"data"`  // Topic data
}
