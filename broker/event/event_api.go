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

import (
	"fmt"

	"github.com/cloustone/sentel/broker/base"
	"github.com/golang/glog"
)

type EventHandler func(e *Event, ctx interface{})

// Publish publish event to event service
func Notify(e *Event) {
	glog.Infof("event '%s' is notified", NameOfEvent(e.Type))
	service := base.GetService(ServiceName).(*eventService)
	service.notify(e)
}

// Subscribe subcribe event from event service
func Subscribe(event uint32, handler EventHandler, ctx interface{}) {
	glog.Infof("service '%s' subscribed event '%s'", ctx.(base.Service).Name(), NameOfEvent(event))
	service := base.GetService(ServiceName).(*eventService)
	service.subscribe(event, handler, ctx)
}

// NameOfEvent return event name
func NameOfEvent(t uint32) string {
	switch t {
	case SessionCreate:
		return "SessionCreate"
	case SessionDestroy:
		return "SessionDestroy"
	case TopicPublish:
		return "TopicPublishe"
	case TopicSubscribe:
		return "TopicSubscribe"
	case TopicUnsubscribe:
		return "TopicUnsubscribe"
	case QutoChange:
		return "QutoChange"
	case SessionResume:
		return "SessionResume"
	case AuthChange:
		return "AuthChange"
	default:
		return "Unknown"
	}
}

// FullNameOfEvent return event information
func FullNameOfEvent(e *Event) string {
	return fmt.Sprintf("Event:%s, broker:%s, clientid:%s", NameOfEvent(e.Type), e.BrokerId, e.ClientId)
}
