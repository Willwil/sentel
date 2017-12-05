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

// sessionManager manage session and subscription and topic tree
// sessionManager receive kafka event for broker cluster's notification
type sessionManager struct {
	base.ServiceBase                    // Session manager is a service that can handle asynchrous event
	eventChan        chan *event.Event  //  Event channel for service
	tree             topicTree          // Global topic subscription tree
	sessions         map[string]Session // All sessions
	mutex            sync.Mutex
}

const (
	ServiceName = "sessionmgr"
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
	glog.Infof("sessionmgr: %s", event.FullNameOfEvent(e))

	switch e.Type {
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
func onEventCallback(e *event.Event, ctx interface{}) {
	service := ctx.(*sessionManager)
	service.eventChan <- e
}

// onEventSessionCreated called when EventSessionCreated event received
func (p *sessionManager) onSessionCreate(e *event.Event) {
	// If session is created on other broker and is retained
	// local queue should be created in local subscription tree
	detail := e.Detail.(*event.SessionCreateDetail)
	if e.BrokerId != base.GetBrokerId() && detail.Persistent {
		// check wethe the session is already exist in subsription tree
		// create virtual session if not exist
		if s, _ := p.findSession(e.ClientId); s == nil {
			if session, _ := newVirtualSession(e.BrokerId, e.ClientId, true); session != nil {
				p.registerSession(session)
			}
		}
	}
}

// onEventSessionDestroyed called when EventSessionDestroyed received
func (p *sessionManager) onSessionDestroy(e *event.Event) {
	session, _ := p.findSession(e.ClientId)
	// For persistent session, change queue's observer
	if session != nil && session.IsPersistent() {
		if q := queue.GetQueue(e.ClientId); q != nil {
			q.RegisterObserver(nil)
		}
		return
	}
	// For transient session, just remove it from session manager
	p.removeSession(e.ClientId)
	queue.DestroyQueue(e.ClientId)
}

// onEventTopicSubscribe called when EventTopicSubscribe received
func (p *sessionManager) onTopicSubscribe(e *event.Event) {
	detail := e.Detail.(*event.TopicSubscribeDetail)

	// Add subscription for persistent subsription from other broker
	if e.BrokerId != base.GetBrokerId() && detail.Persistent {
		glog.Infof("sessionmgr: topic(%s,%s) is subscribed", e.ClientId, detail.Topic)
		if queue := queue.GetQueue(e.ClientId); queue != nil {
			p.tree.addSubscription(&subscription{
				clientId: e.ClientId,
				topic:    detail.Topic,
				qos:      detail.Qos,
				queue:    queue,
			})
		}
	} else if e.BrokerId == base.GetBrokerId() {
		if queue := queue.GetQueue(e.ClientId); queue != nil {
			p.tree.addSubscription(&subscription{
				clientId: e.ClientId,
				topic:    detail.Topic,
				qos:      detail.Qos,
				queue:    queue,
			})
		}
	}
}

// onEventTopicUnsubscribe called when EventTopicUnsubscribe received
func (p *sessionManager) onTopicUnsubscribe(e *event.Event) {
	detail := e.Detail.(*event.TopicUnsubscribeDetail)
	glog.Infof("sessionmgr: topic(%s,%s) is unsubscribed", e.ClientId, detail.Topic)
	if e.BrokerId != base.GetBrokerId() && detail.Persistent {
		p.tree.removeSubscription(e.ClientId, detail.Topic)
	} else if e.BrokerId == base.GetBrokerId() {
		p.tree.removeSubscription(e.ClientId, detail.Topic)
	}
}

// onEventTopicPublish called when EventTopicPublish received
func (p *sessionManager) onTopicPublish(e *event.Event) {
	detail := e.Detail.(*event.TopicPublishDetail)
	glog.Infof("sessionmgr: topic(%s,%s) is published", e.ClientId, detail.Topic)
	p.tree.addMessage(e.ClientId, &base.Message{
		Topic:   detail.Topic,
		Qos:     detail.Qos,
		Payload: detail.Payload,
		Retain:  detail.Retain,
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
