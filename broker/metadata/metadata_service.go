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

	"github.com/cloustone/sentel/broker/broker"
	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"

	"gopkg.in/mgo.v2"
)

// Metaservice manage broker metadata
// Broker's metadata include the following data
// - Global broker cluster data
// - Shadow device
type MetadataService struct {
	core.ServiceBase
	eventChan chan *event.Event
	topicTree TopicTree
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
func (p *MetadataServiceFactory) New(c core.Config, quit chan os.Signal) (core.Service, error) {
	// check mongo db configuration
	hosts, _ := core.GetServiceEndpoint(c, "broker", "mongo")
	timeout := c.MustInt("broker", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	topicTree, err := newTopicTree(c)
	if err != nil {
		return nil, err
	}

	return &MetadataService{
		ServiceBase: core.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		eventChan: make(chan *event.Event),
		topicTree: topicTree,
	}, nil

}

// Name
func (p *MetadataService) Name() string {
	return ServiceName
}

// Start
func (p *MetadataService) Start() error {
	// subscribe envent
	event.Subscribe(event.SessionCreated, onEventCallback, p)
	event.Subscribe(event.SessionDestroyed, onEventCallback, p)
	event.Subscribe(event.TopicSubscribed, onEventCallback, p)
	event.Subscribe(event.TopicUnsubscribed, onEventCallback, p)

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

func (p *MetadataService) handleEvent(e *event.Event) {
	switch e.Type {
	case event.SessionCreated:
		p.onSessionCreated(e)
	case event.SessionDestroyed:
		p.onSessionDestroyed(e)
	case event.TopicSubscribed:
		p.onTopicSubscribe(e)
	case event.TopicUnsubscribed:
		p.onTopicUnsubscribe(e)
	}
}

// onEventCallback will be called when notificaiton come from event service
func onEventCallback(e *event.Event, ctx interface{}) {
	service := ctx.(*MetadataService)
	service.eventChan <- e
}

// onEventSessionCreated called when EventSessionCreated event received
func (p *MetadataService) onSessionCreated(e *event.Event) {
	// Save session info if session is local and retained
	if e.BrokerId == broker.GetId() && e.Persistent {
		// check mongo db configuration
		hosts, _ := core.GetServiceEndpoint(p.Config, "broker", "mongo")
		timeout := p.Config.MustInt("broker", "connect_timeout")
		session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
		if err != nil {
			glog.Errorf("Metadata :%s", err.Error())
			return
		}
		defer session.Close()
		c := session.DB(brokerDatabase).C(sessionCollection)
		err = c.Insert(&Session{
			Id:        e.ClientId,
			CreatedAt: time.Now(),
		})
		if err != nil {
			glog.Error("Metadata:%s", err.Error())
			return
		}
	}

}

// onEventSessionDestroyed called when EventSessionDestroyed received
func (p *MetadataService) onSessionDestroyed(e *event.Event) {
}

// onEventTopicSubscribe called when EventTopicSubscribe received
func (p *MetadataService) onTopicSubscribe(e *event.Event) {
}

// onEventTopicUnsubscribe called when EventTopicUnsubscribe received
func (p *MetadataService) onTopicUnsubscribe(e *event.Event) {
}

// getShadowDeviceStatus return shadow device's status
func (p *MetadataService) getShadowDeviceStatus(clientId string) (*Device, error) {
	return nil, nil
}

// syncShadowDeviceStatus synchronize shadow device's status
func (p *MetadataService) syncShadowDeviceStatus(clientId string, d *Device) error {
	return nil
}

// findSesison return session object by id if existed
func (p *MetadataService) findSession(clientId string) (*Session, error) {
	return p.topicTree.FindSession(clientId)
}

// deleteSession remove session specified by clientId from metadata
func (p *MetadataService) deleteSession(clientId string) error {
	return p.topicTree.DeleteSession(clientId)
}

// regiserSession register session into metadata
func (p *MetadataService) registerSession(s *Session) error {
	return p.topicTree.RegisterSession(s)
}

// addSubscription add a subscription into metadat
func (p *MetadataService) addSubscription(clientId string, topic string, qos uint8) error {
	return p.topicTree.AddSubscription(clientId, topic, qos)
}

// retainSubscription retain the client with topic
func (p *MetadataService) retainSubscription(clientId string, topic string, qos uint8) error {
	return p.topicTree.RetainSubscription(clientId, topic, qos)
}

// RmoeveSubscription remove specified topic from metadata
func (p *MetadataService) removeSubscription(clientId string, topic string) error {
	return p.topicTree.RemoveSubscription(clientId, topic)
}

// deleteMessageWithValidator delete message in metadata with confition
func (p *MetadataService) deleteMessageWithValidator(clientId string, validator func(Message) bool) {
	p.topicTree.DeleteMessageWithValidator(clientId, validator)
}

// deleteMessge delete message specified by idfrom metadata
func (p *MetadataService) deleteMessage(clientId string, mid uint16, direction MessageDirection) error {
	return p.topicTree.DeleteMessage(clientId, mid, direction)
}

// queueMessage save message into metadata
func (p *MetadataService) queueMessage(clientId string, msg *Message) error {
	return p.topicTree.QueueMessage(clientId, *msg)
}

// insertMessage insert specified message into metadata
func (p *MetadataService) insertMessage(clientId string, mid uint16, direction MessageDirection, msg *Message) error {
	return p.topicTree.InsertMessage(clientId, mid, direction, *msg)
}

// releaseMessage release message from metadata
func (p *MetadataService) releaseMessage(clientId string, mid uint16, direction MessageDirection) error {
	return p.topicTree.ReleaseMessage(clientId, mid, direction)
}
