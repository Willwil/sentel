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
	"sync"
	"time"

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/core"

	"gopkg.in/mgo.v2"
)

// Metaservice manage broker metadata
// Broker's metadata include the following data
// - Global broker cluster data
// - Shadow device
type MetadataService struct {
	core.ServiceBase
	eventChan chan *event.Event
}

const (
	ServiceName = "metadata"
)

// MetadataServiceFactory
type MetadataServiceFactory struct{}

// New create metadata service factory
func (p *MetadataServiceFactory) New(c core.Config, quit chan os.Signal) (core.Service, error) {
	// check mongo db configuration
	hosts, _ := core.GetServiceEndpoint(c, "broker", "mongo")
	timeout := c.MustInt("broker", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return &MetadataService{
		ServiceBase: core.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		eventChan: make(chan *event.Event),
	}, nil

}

// Name
func (p *MetadataService) Name() string {
	return ServiceName
}

// Start
func (p *MetadataService) Start() error {
	// subscribe envent
	event.Subscribe(event.EventSessionCreated, onEventCallback, p)
	event.Subscribe(event.EventSessionDestroyed, onEventCallback, p)

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
	p.WaitGroup.Wait()
}

// getShadowDeviceStatus return shadow device's status
func (p *MetadataService) getShadowDeviceStatus(clientId string) (*Device, error) {
	return nil, nil
}

// syncShadowDeviceStatus synchronize shadow device's status
func (p *MetadataService) syncShadowDeviceStatus(clientId string, d *Device) error {
	return nil
}

func (p *MetadataService) handleEvent(e *event.Event) {
	switch e.Type {
	case event.EventSessionCreated:
	case event.EventSessionDestroyed:
	case event.EventTopicPublish:
	case event.EventTopicSubscribe:
	case event.EventTopicUnsubscribe:
	}
}

// onEventCallback will be called when notificaiton come from event service
func onEventCallback(e *event.Event, ctx interface{}) {
	service := ctx.(*MetadataService)
	service.eventChan <- e
}
