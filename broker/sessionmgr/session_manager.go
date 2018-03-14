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
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/broker/queue"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"

	"gopkg.in/mgo.v2"
)

// sessionManager manage session and subscription and topic tree
// sessionManager receive kafka event for broker cluster's notification
type sessionManager struct {
	config    config.Config
	waitgroup sync.WaitGroup
	quitChan  chan interface{}
	eventChan chan event.Event   //  Event channel for service
	tree      topicTree          // Global topic subscription tree
	sessions  map[string]Session // All sessions
	mutex     sync.Mutex
}

const (
	ServiceName = "sessionmgr"
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

	topicTree, err := newTopicTree(c)
	if err != nil {
		return nil, err
	}

	return &sessionManager{
		config:    c,
		waitgroup: sync.WaitGroup{},
		quitChan:  make(chan interface{}),
		eventChan: make(chan event.Event),
		tree:      topicTree,
		sessions:  make(map[string]Session),
		mutex:     sync.Mutex{},
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
	event.Subscribe(event.SessionCreate, onEventCallback, p)
	event.Subscribe(event.SessionDestroy, onEventCallback, p)
	event.Subscribe(event.TopicSubscribe, onEventCallback, p)
	event.Subscribe(event.TopicUnsubscribe, onEventCallback, p)
	event.Subscribe(event.TopicPublish, onEventCallback, p)

	p.waitgroup.Add(1)
	go func(p *sessionManager) {
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
func (p *sessionManager) Stop() {
	p.quitChan <- true
	p.waitgroup.Wait()
	close(p.eventChan)
	close(p.quitChan)
}

func (p *sessionManager) handleEvent(e event.Event) {
	glog.Infof("sessionmgr: %s", event.FullNameOfEvent(e))

	switch e.GetType() {
	case event.SessionCreate:
		p.onSessionCreate(e)
	case event.SessionDestroy:
		p.onSessionDestroy(e)
	case event.TopicSubscribe:
		p.onTopicSubscribe(e)
	case event.TopicUnsubscribe:
		p.onTopicUnsubscribe(e)
	case event.TopicPublish:
		p.onTopicPublish(e)
	}
}

// onEventCallback will be called when notificaiton come from event service
func onEventCallback(e event.Event, ctx interface{}) {
	service := ctx.(*sessionManager)
	service.eventChan <- e
}

// onEventSessionCreated called when EventSessionCreated event received
func (p *sessionManager) onSessionCreate(e event.Event) {
	// If session is created on other broker and is retained
	// local queue should be created in local subscription tree
	if e.GetBrokerId() != base.GetBrokerId() {
		if sce, ok := e.(*event.SessionCreateEvent); ok && sce.Persistent {
			glog.Infof("sessionmgr receive session create notification from cluster broker '%s'", e.GetBrokerId())
			// check wethe the session is already exist in subsription tree
			// create virtual session if not exist
			if s, _ := p.findSession(e.GetClientId()); s == nil {
				if session, _ := newVirtualSession(e.GetBrokerId(), e.GetClientId(), true); session != nil {
					p.registerSession(session)
				}
			}
		}
	}
}

// onEventSessionDestroyed called when EventSessionDestroyed received
func (p *sessionManager) onSessionDestroy(e event.Event) {
	clientId := e.GetClientId()
	session, _ := p.findSession(clientId)
	// For persistent session, change queue's observer
	if session != nil && session.IsPersistent() {
		if q := queue.GetQueue(clientId); q != nil {
			q.RegisterObserver(nil)
		}
		return
	}
	// For transient session, just remove it from session manager
	p.removeSession(clientId)
	queue.DestroyQueue(clientId)
}

// onEventTopicSubscribe called when EventTopicSubscribe received
func (p *sessionManager) onTopicSubscribe(e event.Event) {
	t := e.(*event.TopicSubscribeEvent)

	// Add subscription for persistent subsription from other broker
	if t.BrokerID != base.GetBrokerId() && t.Persistent {
		glog.Infof("sessionmgr receive topic('%s') subscribe notification from cluster broker'%s'", t.Topic, t.ClientID)
		if queue := queue.GetQueue(t.ClientID); queue != nil {
			p.tree.addSubscription(&subscription{
				clientID: t.ClientID,
				topic:    t.Topic,
				qos:      t.Qos,
				queue:    queue,
			})
		}
	} else if t.BrokerID == base.GetBrokerId() {
		if queue := queue.GetQueue(t.ClientID); queue != nil {
			p.tree.addSubscription(&subscription{
				clientID: t.ClientID,
				topic:    t.Topic,
				qos:      t.Qos,
				queue:    queue,
			})
		}
	}
}

// onEventTopicUnsubscribe called when EventTopicUnsubscribe received
func (p *sessionManager) onTopicUnsubscribe(e event.Event) {
	t := e.(*event.TopicUnsubscribeEvent)
	glog.Infof("sessionmgr: topic(%s,%s) is unsubscribed", t.ClientID, t.Topic)
	if t.BrokerID != base.GetBrokerId() && t.Persistent {
		p.tree.removeSubscription(t.ClientID, t.Topic)
	} else if t.BrokerID == base.GetBrokerId() {
		p.tree.removeSubscription(t.ClientID, t.Topic)
	}
}

// onEventTopicPublish called when EventTopicPublish received
func (p *sessionManager) onTopicPublish(e event.Event) {
	t := e.(*event.TopicPublishEvent)
	glog.Infof("sessionmgr: topic(%s,%s) is published", t.ClientID, t.Topic)
	p.tree.addMessage(t.ClientID, &base.Message{
		Topic:   t.Topic,
		Qos:     t.Qos,
		Payload: t.Payload,
		Retain:  t.Retain,
	})
}

// findSesison return session object by id if existed
func (p *sessionManager) findSession(clientId string) (Session, error) {
	glog.Infof("sessionmgr: find session '%s'", clientId)
	if _, ok := p.sessions[clientId]; !ok {
		return nil, errors.New("sessionmgr:Session id does not exist")
	}

	return p.sessions[clientId], nil

}

// deleteSession remove session specified by clientId from metadata
func (p *sessionManager) removeSession(clientId string) error {
	glog.Infof("sessionmgr: remove session for '%s'", clientId)
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.sessions[clientId]; !ok {
		return fmt.Errorf("Session '%s' does not exist", clientId)
	}
	delete(p.sessions, clientId)
	return nil
}

// regiserSession register session into metadata
func (p *sessionManager) registerSession(s Session) error {
	glog.Infof("sessionmgr: register session for '%s'", s.Id())
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.sessions[s.Id()]; ok {
		return errors.New("Session id already exists")
	}

	p.sessions[s.Id()] = s
	return nil
}

// deleteMessageWithValidator delete message in metadata with confition
func (p *sessionManager) deleteMessageWithValidator(clientId string, validator func(*base.Message) bool) {
	glog.Infof("sessionmgr: delete message for %s", clientId)
	p.tree.deleteMessageWithValidator(clientId, validator)
}

// deleteMessge delete message specified by idfrom metadata
func (p *sessionManager) deleteMessage(clientId string, mid uint16, direction uint8) {
	p.deleteMessageWithValidator(clientId, func(msg *base.Message) bool {
		return msg.PacketId == mid && msg.Direction == direction
	})
}

// getTopics return all topic in the broker
func (p *sessionManager) getTopics() []*Topic {
	return nil
}

// getClientTopicss return client's subscribed topics
func (p *sessionManager) getClientTopics(clientId string) []*Topic {
	return nil
}

// getTopic return specified topic subscription info
func (p *sessionManager) getTopicSubscription(topic string) []*Subscription {
	return nil
}

// getSubscriptions return all subscriptions in the broker
func (p *sessionManager) getSubscriptions() []*Subscription {
	return nil
}

// getClientSubscriptions return client's subscription
func (p *sessionManager) getClientSubscriptions(clientId string) []*Subscription {
	return nil
}

//  getSession return all sessions in the broker
func (p *sessionManager) getSessions() []*Session {
	return nil
}

// getSessionInfo return client's session detail information
func (p *sessionManager) getSessionInfo(clientId string) *Session {
	return nil
}

// getClients return all clients in broker cluster
func (p *sessionManager) getClients() []*Client {
	return nil
}

// getClient return client's detail information
func (p *sessionManager) getClient(clientId string) *Client {
	return nil
}
