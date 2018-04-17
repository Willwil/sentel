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

package metadata

import (
	"sync"
	"time"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"

	"gopkg.in/mgo.v2"
)

// Metaservice manage broker metadata
// Broker's metadata include the following data
// - Global broker cluster data
// - Shadow device
type metadataService struct {
	config    config.Config
	waitgroup sync.WaitGroup
	quitChan  chan interface{}
	eventChan chan event.Event
}

const (
	ServiceName       = "metadata"
	brokerDatabase    = "broker"
	sessionCollection = "sessions"
	brokerCollection  = "brokers"
)

type ServiceFactory struct{}

// New create metadata service factory
func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	// check mongo db configuration
	hosts := c.MustString("mongo")
	timeout := c.MustInt("connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return &metadataService{
		config:    c,
		waitgroup: sync.WaitGroup{},
		quitChan:  make(chan interface{}),
		eventChan: make(chan event.Event),
	}, nil

}

// Name
func (p *metadataService) Name() string {
	return ServiceName
}

func (p *metadataService) Initialize() error { return nil }

// Start
func (p *metadataService) Start() error {
	// subscribe envent
	event.Subscribe(event.SessionCreate, onEventCallback, p)
	event.Subscribe(event.SessionDestroy, onEventCallback, p)
	event.Subscribe(event.TopicSubscribe, onEventCallback, p)
	event.Subscribe(event.TopicUnsubscribe, onEventCallback, p)

	p.waitgroup.Add(1)
	go func(p *metadataService) {
		defer p.waitgroup.Done()
		for {
			select {
			case e := <-p.eventChan:
				p.handleEvent(e)
			case <-p.quitChan:
				return
			}
		}
	}(p)
	return nil
}

// Stop
func (p *metadataService) Stop() {
	p.quitChan <- true
	p.waitgroup.Wait()
	close(p.quitChan)
	close(p.eventChan)
}

func (p *metadataService) handleEvent(e event.Event) {
	switch e.GetType() {
	case event.SessionCreate:
		p.onSessionCreated(e)
	case event.SessionDestroy:
		p.onSessionDestroyed(e)
	case event.TopicSubscribe:
		p.onTopicSubscribe(e)
	case event.TopicUnsubscribe:
		p.onTopicUnsubscribe(e)
	}
}

// onEventCallback will be called when notificaiton come from event service
func onEventCallback(e event.Event, ctx interface{}) {
	service := ctx.(*metadataService)
	service.eventChan <- e
}

// onEventSessionCreated called when EventSessionCreated event received
func (p *metadataService) onSessionCreated(e event.Event) {
}

// onEventSessionDestroyed called when EventSessionDestroyed received
func (p *metadataService) onSessionDestroyed(e event.Event) {
}

// onEventTopicSubscribe called when EventTopicSubscribe received
func (p *metadataService) onTopicSubscribe(e event.Event) {
}

// onEventTopicUnsubscribe called when EventTopicUnsubscribe received
func (p *metadataService) onTopicUnsubscribe(e event.Event) {
}

// getShadowDeviceStatus return shadow device's status
func (p *metadataService) getShadowDeviceStatus(clientID string) (*Device, error) {
	return nil, nil
}

// syncShadowDeviceStatus synchronize shadow device's status
func (p *metadataService) syncShadowDeviceStatus(clientID string, d *Device) error {
	return nil
}

// deleteMessageWithValidator delete message in metadata with confition
func (p *metadataService) deleteMessageWithValidator(clientID string, validator func(*base.Message) bool) {
}

// deleteMessge delete message specified by idfrom metadata
func (p *metadataService) deleteMessage(clientID string, pid uint16, direction uint8) {
	p.deleteMessageWithValidator(clientID, func(msg *base.Message) bool {
		return msg.PacketId == pid && msg.Direction == direction
	})
}

// findMessage retrieve message with packet id
func (p *metadataService) findMessage(clientid string, pid uint16) *base.Message {
	return nil
}
