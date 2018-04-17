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

import "encoding/json"

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

type Event interface {
	SetBrokerId(brokerId string)
	SetType(eventType uint32)
	SetClientId(clientID string)
	GetBrokerId() string
	GetType() uint32
	GetClientId() string
	//Serialize(*CodecOption) ([]byte, error)
}

// RawEvent is event pubished to message server
type rawEvent struct {
	Type    uint32          `json:"type"` // Event type
	Payload json.RawMessage `json:"payload"`
}
