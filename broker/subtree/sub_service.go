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

package subtree

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/broker/broker"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"

	"gopkg.in/mgo.v2"
)

// Metaservice manage broker metadata
// Broker's metadata include the following data
// - Global broker cluster data
// - Shadow device
type SubTreeService struct {
	base.ServiceBase
	eventChan chan *broker.Event
	topicTree TopicTree
}

const (
	ServiceName       = "subtree"
	brokerDatabase    = "broker"
	sessionCollection = "sessions"
	brokerCollection  = "brokers"
)

// SubTreeServiceFactory
type SubServiceFactory struct{}

// New create metadata service factory
func (p *SubServiceFactory) New(c core.Config, quit chan os.Signal) (base.Service, error) {
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

	return &SubTreeService{
		ServiceBase: base.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		eventChan: make(chan *broker.Event),
		topicTree: topicTree,
	}, nil

}

// Name
func (p *SubTreeService) Name() string {
	return ServiceName
}

// Start
func (p *SubTreeService) Start() error {
	// subscribe envent
	broker.Subscribe(broker.SessionCreated, onEventCallback, p)
	broker.Subscribe(broker.SessionDestroyed, onEventCallback, p)
	broker.Subscribe(broker.TopicSubscribed, onEventCallback, p)
	broker.Subscribe(broker.TopicUnsubscribed, onEventCallback, p)
	broker.Subscribe(broker.TopicPublished, onEventCallback, p)

	go func(p *SubTreeService) {
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
func (p *SubTreeService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
	close(p.eventChan)
}

func (p *SubTreeService) handleEvent(e *broker.Event) {
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
	service := ctx.(*SubTreeService)
	service.eventChan <- e
}

// onEventSessionCreated called when EventSessionCreated event received
func (p *SubTreeService) onSessionCreated(e *broker.Event) {
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
func (p *SubTreeService) onSessionDestroyed(e *broker.Event) {
}

// onEventTopicSubscribe called when EventTopicSubscribe received
func (p *SubTreeService) onTopicSubscribe(e *broker.Event) {
}

// onEventTopicUnsubscribe called when EventTopicUnsubscribe received
func (p *SubTreeService) onTopicUnsubscribe(e *broker.Event) {
}

// findSesison return session object by id if existed
func (p *SubTreeService) findSession(clientId string) (*Session, error) {
	return p.topicTree.FindSession(clientId)
}

// deleteSession remove session specified by clientId from metadata
func (p *SubTreeService) deleteSession(clientId string) error {
	return p.topicTree.DeleteSession(clientId)
}

// regiserSession register session into metadata
func (p *SubTreeService) registerSession(s *Session) error {
	return p.topicTree.RegisterSession(s)
}

// addSubscription add a subscription into metadat
func (p *SubTreeService) addSubscription(clientId string, topic string, qos uint8) error {
	return p.topicTree.AddSubscription(clientId, topic, qos)
}

// retainSubscription retain the client with topic
func (p *SubTreeService) retainSubscription(clientId string, topic string, qos uint8) error {
	return p.topicTree.RetainSubscription(clientId, topic, qos)
}

// RmoeveSubscription remove specified topic from metadata
func (p *SubTreeService) removeSubscription(clientId string, topic string) error {
	return p.topicTree.RemoveSubscription(clientId, topic)
}

// deleteMessageWithValidator delete message in metadata with confition
func (p *SubTreeService) deleteMessageWithValidator(clientId string, validator func(Message) bool) {
	p.topicTree.DeleteMessageWithValidator(clientId, validator)
}

// deleteMessge delete message specified by idfrom metadata
func (p *SubTreeService) deleteMessage(clientId string, mid uint16, direction MessageDirection) error {
	return p.topicTree.DeleteMessage(clientId, mid, direction)
}

// queueMessage save message into metadata
func (p *SubTreeService) queueMessage(clientId string, msg *Message) error {
	return p.topicTree.QueueMessage(clientId, *msg)
}

// insertMessage insert specified message into metadata
func (p *SubTreeService) insertMessage(clientId string, mid uint16, direction MessageDirection, msg *Message) error {
	return p.topicTree.InsertMessage(clientId, mid, direction, *msg)
}

// releaseMessage release message from metadata
func (p *SubTreeService) releaseMessage(clientId string, mid uint16, direction MessageDirection) error {
	return p.topicTree.ReleaseMessage(clientId, mid, direction)
}
