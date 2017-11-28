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

package sessionmgr

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/broker/queue"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"

	"gopkg.in/mgo.v2"
)

type sessionManager struct {
	base.ServiceBase
	eventChan chan *event.Event
	topicTree TopicTree
}

const (
	ServiceName       = "subtree"
	brokerDatabase    = "broker"
	sessionCollection = "sessions"
	brokerCollection  = "brokers"
)

// New create metadata service factory
func New(c core.Config, quit chan os.Signal) (base.Service, error) {
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

	return &sessionManager{
		ServiceBase: base.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		eventChan: make(chan *event.Event),
		topicTree: topicTree,
	}, nil

}

// Name
func (p *sessionManager) Name() string {
	return ServiceName
}

func (p *sessionManager) Initialize() error { return nil }

// Start
func (p *sessionManager) Start() error {
	// subscribe envent
	event.Subscribe(event.SessionCreated, onEventCallback, p)
	event.Subscribe(event.SessionDestroyed, onEventCallback, p)
	event.Subscribe(event.TopicSubscribed, onEventCallback, p)
	event.Subscribe(event.TopicUnsubscribed, onEventCallback, p)
	event.Subscribe(event.TopicPublished, onEventCallback, p)

	go func(p *sessionManager) {
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
func (p *sessionManager) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
	close(p.eventChan)
}

func (p *sessionManager) handleEvent(e *event.Event) {
	switch e.Type {
	case event.SessionCreated:
		p.onSessionCreate(e)
	case event.SessionDestroyed:
		p.onSessionDestroy(e)
	case event.TopicSubscribed:
		p.onTopicSubscribe(e)
	case event.TopicUnsubscribed:
		p.onTopicUnsubscribe(e)
	case event.TopicPublished:
		p.onTopicPublish(e)
	}
}

// onEventCallback will be called when notificaiton come from event service
func onEventCallback(e *event.Event, ctx interface{}) {
	service := ctx.(*sessionManager)
	service.eventChan <- e
}

// onEventSessionCreated called when EventSessionCreated event received
func (p *sessionManager) onSessionCreate(e *event.Event) {
	// If session is created on other broker and is retained
	// local queue should be created in local subscription tree
	if e.BrokerId != base.GetBrokerId() && e.Persistent {

		// check wethe the session is already exist in subsription tree
		// return simpily if session is already retained
		session, err := p.topicTree.FindSession(e.ClientId)
		if err == nil && session.Retain == true {
			return
		}
		// create queue if not exist
		_, err = queue.NewQueue(e.ClientId, true)
		if err != nil {
			glog.Fatalf("broker: Failed to create queue for client '%s'", e.ClientId)
		}
	}

}

// onEventSessionDestroyed called when EventSessionDestroyed received
func (p *sessionManager) onSessionDestroy(e *event.Event) {
	glog.Infof("subtree: session(%s) is destroyed", e.ClientId)
}

// onEventTopicSubscribe called when EventTopicSubscribe received
func (p *sessionManager) onTopicSubscribe(e *event.Event) {
	glog.Infof("subtree: topic(%s,%s) is subscribed", e.ClientId, e.Topic)
	queue := queue.GetQueue(e.ClientId)
	if queue != nil {
		glog.Fatalf("broker: Can not get queue for client '%s'", e.ClientId)
		return
	}
	p.topicTree.AddSubscription(e.ClientId, e.Topic, e.Qos, queue)
}

// onEventTopicUnsubscribe called when EventTopicUnsubscribe received
func (p *sessionManager) onTopicUnsubscribe(e *event.Event) {
	glog.Infof("subtree: topic(%s,%s) is unsubscribed", e.ClientId, e.Topic)
	p.topicTree.RemoveSubscription(e.ClientId, e.Topic)
}

// onEventTopicPublish called when EventTopicPublish received
func (p *sessionManager) onTopicPublish(e *event.Event) {
	glog.Infof("substree: topic(%s,%s) is published", e.ClientId, e.Topic)
	p.topicTree.AddTopic(e.ClientId, e.Topic, e.Data)
}

// findSesison return session object by id if existed
func (p *sessionManager) findSession(clientId string) (*Session, error) {
	return p.topicTree.FindSession(clientId)
}

// deleteSession remove session specified by clientId from metadata
func (p *sessionManager) deleteSession(clientId string) error {
	return p.topicTree.DeleteSession(clientId)
}

// regiserSession register session into metadata
func (p *sessionManager) registerSession(s *Session) error {
	return p.topicTree.RegisterSession(s)
}

// addSubscription add a subscription into metadat
func (p *sessionManager) addSubscription(clientId string, topic string, qos uint8, q queue.Queue) error {
	return p.topicTree.AddSubscription(clientId, topic, qos, q)
}

// retainSubscription retain the client with topic
func (p *sessionManager) retainSubscription(clientId string, topic string, qos uint8) error {
	return p.topicTree.RetainSubscription(clientId, topic, qos)
}

// RmoeveSubscription remove specified topic from metadata
func (p *sessionManager) removeSubscription(clientId string, topic string) error {
	return p.topicTree.RemoveSubscription(clientId, topic)
}

// deleteMessageWithValidator delete message in metadata with confition
func (p *sessionManager) deleteMessageWithValidator(clientId string, validator func(Message) bool) {
	p.topicTree.DeleteMessageWithValidator(clientId, validator)
}

// deleteMessge delete message specified by idfrom metadata
func (p *sessionManager) deleteMessage(clientId string, mid uint16, direction MessageDirection) error {
	return p.topicTree.DeleteMessage(clientId, mid, direction)
}

// queueMessage save message into metadata
func (p *sessionManager) queueMessage(clientId string, msg *Message) error {
	return p.topicTree.QueueMessage(clientId, *msg)
}

// insertMessage insert specified message into metadata
func (p *sessionManager) insertMessage(clientId string, mid uint16, direction MessageDirection, msg *Message) error {
	return p.topicTree.InsertMessage(clientId, mid, direction, *msg)
}

// releaseMessage release message from metadata
func (p *sessionManager) releaseMessage(clientId string, mid uint16, direction MessageDirection) error {
	return p.topicTree.ReleaseMessage(clientId, mid, direction)
}
