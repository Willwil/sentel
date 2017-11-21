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

import "github.com/cloustone/sentel/broker/base"

type EventHandler func(ctx interface{}, e *Event)

const (
	EventSessionCreated   = 1
	EventSessionDestroyed = 2
	EventTopicPublish     = 3
	EventTopicSubscribe   = 4
	EventTopicUnsubscribe = 5
)

type Event struct {
	Type       uint8  // Event type
	ClientId   string // Client identifier where event come from
	SessionId  string // Current Session identifier
	Topic      string // Topic
	Data       []byte // Topic data
	Persistent bool   // Whether the session is persistent
}

// Publish publish event to event service
func Publish(e *Event) {
	service := base.GetService(EventServiceName).(*EventService)
	service.publish(e)
}

// Subscribe subcribe event from event service
func Subscribe(t uint8, handler EventHandler, ctx interface{}) {
	service := base.GetService(EventServiceName).(*EventService)
	service.subscribe(t, handler, ctx)
}
