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
	"github.com/cloustone/sentel/broker/broker"
	"github.com/cloustone/sentel/core"

	"gopkg.in/mgo.v2"
)

// Metaservice manage broker metadata
// Broker's metadata include the following data
// - Global broker cluster data
// - Shadow device
type MetadataService struct {
	base.ServiceBase
	eventChan chan *broker.Event
}

const (
	ServiceName       = "metadata"
	brokerDatabase    = "broker"
	sessionCollection = "sessions"
	brokerCollection  = "brokers"
)

// MetadataServiceFactory
type MetadataServiceFactory struct{}

// New create metadata service factory
func (p *MetadataServiceFactory) New(c core.Config, quit chan os.Signal) (base.Service, error) {
	// check mongo db configuration
	hosts, _ := core.GetServiceEndpoint(c, "broker", "mongo")
	timeout := c.MustInt("broker", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return &MetadataService{
		ServiceBase: base.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		eventChan: make(chan *broker.Event),
	}, nil

}

// Name
func (p *MetadataService) Name() string {
	return ServiceName
}

// Start
func (p *MetadataService) Start() error {
	// subscribe envent
	broker.Subscribe(broker.SessionCreated, onEventCallback, p)
	broker.Subscribe(broker.SessionDestroyed, onEventCallback, p)
	broker.Subscribe(broker.TopicSubscribed, onEventCallback, p)
	broker.Subscribe(broker.TopicUnsubscribed, onEventCallback, p)

	go func(p *MetadataService) {
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
func (p *MetadataService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
	close(p.eventChan)
}

func (p *MetadataService) handleEvent(e *broker.Event) {
	switch e.Type {
	case broker.SessionCreated:
		p.onSessionCreated(e)
	case broker.SessionDestroyed:
		p.onSessionDestroyed(e)
	case broker.TopicSubscribed:
		p.onTopicSubscribe(e)
	case broker.TopicUnsubscribed:
		p.onTopicUnsubscribe(e)
	}
}

// onEventCallback will be called when notificaiton come from event service
func onEventCallback(e *broker.Event, ctx interface{}) {
	service := ctx.(*MetadataService)
	service.eventChan <- e
}

// onEventSessionCreated called when EventSessionCreated event received
func (p *MetadataService) onSessionCreated(e *broker.Event) {
	// Save session info if session is local and retained
	if e.BrokerId == broker.GetId() && e.Persistent {
		// check mongo db configuration
	}

}

// onEventSessionDestroyed called when EventSessionDestroyed received
func (p *MetadataService) onSessionDestroyed(e *broker.Event) {
}

// onEventTopicSubscribe called when EventTopicSubscribe received
func (p *MetadataService) onTopicSubscribe(e *broker.Event) {
}

// onEventTopicUnsubscribe called when EventTopicUnsubscribe received
func (p *MetadataService) onTopicUnsubscribe(e *broker.Event) {
}

// getShadowDeviceStatus return shadow device's status
func (p *MetadataService) getShadowDeviceStatus(clientId string) (*Device, error) {
	return nil, nil
}

// syncShadowDeviceStatus synchronize shadow device's status
func (p *MetadataService) syncShadowDeviceStatus(clientId string, d *Device) error {
	return nil
}
