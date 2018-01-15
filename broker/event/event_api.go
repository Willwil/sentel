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
	"encoding/json"
	"fmt"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"
)

// EventHandler event handler
type EventHandler func(e *Event, ctx interface{})

// Notify publish event to event service
func Notify(event uint32, clientId string, detail interface{}) {
	glog.Infof("event '%s' is notified", NameOfEvent(event))
	e := Event{
		EventHeader: EventHeader{
			Type:     event,
			ClientId: clientId,
			BrokerId: base.GetBrokerId(),
		},
		Detail: detail,
	}
	service := base.GetService(ServiceName).(*eventService)
	service.notify(&e)
}

// Subscribe subcribe event from event service
func Subscribe(event uint32, handler EventHandler, ctx interface{}) {
	glog.Infof("service '%s' subscribed event '%s'",
		ctx.(service.Service).Name(), NameOfEvent(event))
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
		return "TopicPublish"
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

// FromRawEvent unmarshal event from raw event
func FromRawEvent(value []byte) (*Event, error) {
	re := RawEvent{}
	if err := json.Unmarshal(value, &re); err != nil {
		glog.Errorf("conductor unmarshal event common failed:%s", err.Error())
		return nil, err
	}
	e := &Event{}
	if err := json.Unmarshal(re.Header, &e.EventHeader); err != nil {
		glog.Errorf("conductor unmarshal event common failed:%s", err.Error())
		return nil, err
	}
	switch e.Type {
	case SessionCreate:
		r := SessionCreateDetail{}
		if err := json.Unmarshal(re.Payload, &r); err != nil {
			return nil, fmt.Errorf("event unmarshal failed, %s", err.Error())
		}
		e.Detail = r
	case SessionDestroy:
		r := SessionDestroyDetail{}
		if err := json.Unmarshal(re.Payload, &r); err != nil {
			return nil, fmt.Errorf("event unmarshal failed, %s", err.Error())
		}
		e.Detail = r
	case TopicSubscribe:
		r := TopicSubscribeDetail{}
		if err := json.Unmarshal(re.Payload, &r); err != nil {
			return nil, fmt.Errorf("event unmarshal failed, %s", err.Error())
		}
		e.Detail = r
	case TopicUnsubscribe:
		r := TopicUnsubscribeDetail{}
		if err := json.Unmarshal(re.Payload, &r); err != nil {
			return nil, fmt.Errorf("event unmarshal failed, %s", err.Error())
		}
		e.Detail = r
	case TopicPublish:
		r := TopicPublishDetail{}
		if err := json.Unmarshal(re.Payload, &r); err != nil {
			return nil, fmt.Errorf("event unmarshal failed, %s", err.Error())
		}
		e.Detail = r
	default:
		return nil, fmt.Errorf("invalid event type '%d'", e.Type)
	}
	return e, nil
}
