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
	SessionCreated    = 1
	SessionDestroyed  = 2
	SessionResumed    = 3
	TopicPublished    = 4
	TopicSubscribed   = 5
	TopicUnsubscribed = 6
	QutoChanged       = 7
)

const (
	EventBusName = "broker/event"
)

type Event struct {
	core.TopicBase
	BrokerId   string `json:"brokerId"`   // Broker identifier where event come from
	Type       uint8  `json:"type"`       // Event type
	ClientId   string `json:"clientId"`   // Client identifier where event come from
	Topic      string `json:"topic"`      // Topic
	Data       []byte `json:"data"`       // Topic data
	Persistent bool   `json:"persistent"` // Whether the session is persistent

	// Quto
	QutoId string `json:"qutoId"` // Qutotation identifier
}
