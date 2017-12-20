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
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/common"

	"gopkg.in/mgo.v2"
)

// Metaservice manage broker metadata
// Broker's metadata include the following data
// - Global broker cluster data
// - Shadow device
type metadataService struct {
	base.ServiceBase
	eventChan chan *event.Event
}

const (
	ServiceName       = "metadata"
	brokerDatabase    = "broker"
	sessionCollection = "sessions"
	brokerCollection  = "brokers"
)

// New create metadata service factory
func New(c com.Config, quit chan os.Signal) (base.Service, error) {
	// check mongo db configuration
	hosts := c.MustString("broker", "mongo")
	timeout := c.MustInt("broker", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return &metadataService{
		ServiceBase: base.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		eventChan: make(chan *event.Event),
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

	go func(p *metadataService) {
		for {
			select {
			case e := <-p.eventChan:
				p.handleEvent(e)
			case <-p.Quit:
				return
			}
		}
	}(p)
	return nil
}

// Stop
func (p *metadataService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
	close(p.eventChan)
}

func (p *metadataService) handleEvent(e *event.Event) {
	switch e.Common.Type {
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
func onEventCallback(e *event.Event, ctx interface{}) {
	service := ctx.(*metadataService)
	service.eventChan <- e
}

// onEventSessionCreated called when EventSessionCreated event received
func (p *metadataService) onSessionCreated(e *event.Event) {
}

// onEventSessionDestroyed called when EventSessionDestroyed received
func (p *metadataService) onSessionDestroyed(e *event.Event) {
}

// onEventTopicSubscribe called when EventTopicSubscribe received
func (p *metadataService) onTopicSubscribe(e *event.Event) {
}

// onEventTopicUnsubscribe called when EventTopicUnsubscribe received
func (p *metadataService) onTopicUnsubscribe(e *event.Event) {
}

// getShadowDeviceStatus return shadow device's status
func (p *metadataService) getShadowDeviceStatus(clientId string) (*Device, error) {
	return nil, nil
}

// syncShadowDeviceStatus synchronize shadow device's status
func (p *metadataService) syncShadowDeviceStatus(clientId string, d *Device) error {
	return nil
}

// deleteMessageWithValidator delete message in metadata with confition
func (p *metadataService) deleteMessageWithValidator(clientId string, validator func(*base.Message) bool) {
}

// deleteMessge delete message specified by idfrom metadata
func (p *metadataService) deleteMessage(clientId string, pid uint16, direction uint8) {
	p.deleteMessageWithValidator(clientId, func(msg *base.Message) bool {
		return msg.PacketId == pid && msg.Direction == direction
	})
}

// findMessage retrieve message with packet id
func (p *metadataService) findMessage(clientid string, pid uint16) *base.Message {
	return nil
}
