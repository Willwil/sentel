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
	EventSessionCreated   = 1
	EventSessionDestroyed = 2
	EventTopicPublish     = 3
	EventTopicSubscribe   = 4
	EventTopicUnsubscribe = 5
)

const (
	EventBusName = "broker/event"
)

type Event struct {
	core.TopicBase
	BrokerId   string `json:"brokerId"`   // Broker identifier where event come from
	Type       uint8  `json:"type"`       // Event type
	ClientId   string `json:"clientId"`   // Client identifier where event come from
	SessionId  string `json:"sessionId"`  // Current Session identifier
	Topic      string `json:"topic"`      // Topic
	Data       []byte `json:"data"`       // Topic data
	Persistent bool   `json:"persistent"` // Whether the session is persistent
}
