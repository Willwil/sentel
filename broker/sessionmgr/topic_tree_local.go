//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain
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
	"strings"

	"github.com/cloustone/sentel/broker/queue"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

type subLeaf struct {
	qos uint8
}

type subNode struct {
	level     string
	children  map[string]*subNode
	subs      map[string]*subLeaf
	retainMsg *Message
}

type localTopicTree struct {
	config   core.Config
	sessions map[string]*Session
	root     subNode
}

// FindSession find session by id
func (p *localTopicTree) findSession(id string) (*Session, error) {
	if _, ok := p.sessions[id]; !ok {
		return nil, errors.New("Session id does not exist")
	}

	return p.sessions[id], nil
}

// DeleteSession delete session by id
func (p *localTopicTree) deleteSession(id string) error {
	if _, ok := p.sessions[id]; !ok {
		return errors.New("Session id does not exist")
	}
	delete(p.sessions, id)
	return nil
}

// UpdateSession update session
func (p *localTopicTree) updateSession(s *Session) error {
	if _, ok := p.sessions[s.Id]; !ok {
		return errors.New("Session id does not exist")
	}
	p.sessions[s.Id] = s
	return nil
}

// RegisterSession register new session
func (p *localTopicTree) registerSession(s *Session) error {
	if _, ok := p.sessions[s.Id]; ok {
		return errors.New("Session id already exists")
	}

	glog.Infof("RegisterSession: id is %s", s.Id)
	p.sessions[s.Id] = s
	return nil
}

func (p *localTopicTree) findNode(node *subNode, lev string) *subNode {
	for k, v := range node.children {
		if k == lev {
			return v
		}
	}

	return nil
}

func (p *localTopicTree) addNode(node *subNode, lev string) *subNode {
	for k, v := range node.children {
		if k == lev {
			return v
		}
	}

	tmp := &subNode{
		level:    lev,
		children: make(map[string]*subNode),
		subs:     make(map[string]*subLeaf),
	}

	node.children[lev] = tmp

	return tmp
}

// Subscription
func (p *localTopicTree) addSubscription(sessionid string, topic string, qos uint8, q queue.Queue) error {
	glog.Infof("AddSubscription: sessionid is %s, topic is %s, qos is %d", sessionid, topic, qos)
	node := &p.root
	s := strings.Split(topic, "/")
	for _, level := range s {
		glog.Infof("AddSubscription: level is %s", level)
		node = p.addNode(node, level)
	}

	glog.Infof("AddSubscription: session id is %s", sessionid)
	node.subs[sessionid] = &subLeaf{
		qos: qos,
	}

	return nil
}

// RetainSubscription process RETAIN flag；
func (p *localTopicTree) retainSubscription(sessionid string, topic string, qos uint8) error {
	return nil
}

func (p *localTopicTree) removeSubscription(sessionid string, topic string) error {
	node := &p.root
	s := strings.Split(topic, "/")
	for _, level := range s {
		node = p.findNode(node, level)
		if node == nil {
			return nil
		}
	}

	if _, ok := node.subs[sessionid]; ok {
		delete(node.subs, sessionid)
	}

	return nil
}

// Message Management
func (p *localTopicTree) findMessage(clientid string, mid uint16) (*Message, error) {
	return nil, nil
}

func (p *localTopicTree) storeMessage(clientid string, msg Message) error {
	return nil
}

func (p *localTopicTree) deleteMessageWithValidator(clientid string, validator func(msg Message) bool) {

}

func (p *localTopicTree) deleteMessage(clientid string, mid uint16, direction MessageDirection) error {
	return nil
}

func (p *localTopicTree) subProcess(clientid string, msg *Message, node *subNode, setRetain bool) error {
	/*
		if msg.Retain && setRetain {
			node.retainMsg = msg
		}

		for k, v := range node.subs {
			glog.Infof("subProcess: session id is %s", k)
			s, ok := p.sessions[k]
			if !ok {
				glog.Errorf("subProcess: sessions is nil")
				continue
			}
			// if s.id == clientid {
			// 	continue
			// }

			//s.sendPublish(v.qos, msg.Qos, msg.Topic, msg.Payload)
		}
	*/
	return nil
}

func (p *localTopicTree) subSearch(clientid string, msg *Message, node *subNode, levels []string, setRetain bool) error {
	for k, v := range node.children {
		sr := setRetain
		if len(levels) != 0 && (k == levels[0] || k == "+") {
			if k == "+" {
				sr = false
			}
			ss := levels[1:]
			p.subSearch(clientid, msg, v, ss, sr)
			if len(ss) == 0 {
				p.subProcess(clientid, msg, v, sr)
			}
		} else if k == "#" && len(v.children) > 0 {
			p.subProcess(clientid, msg, v, sr)
		}
	}
	return nil
}

func (p *localTopicTree) queueMessage(clientid string, msg Message) error {
	glog.Infof("QueueMessage: Message Topic is %s", msg.Topic)
	s := strings.Split(msg.Topic, "/")

	if msg.Retain {
		/* We have a message that needs to be retained, so ensure that the subscription
		 * tree for its topic exists.
		 */
		node := &p.root
		for _, level := range s {
			node = p.addNode(node, level)
		}
	}

	return p.subSearch(clientid, &msg, &p.root, s, true)
}

func (p *localTopicTree) getMessageTotalCount(clientid string) int {
	return 0
}

func (p *localTopicTree) insertMessage(clientid string, mid uint16, direction MessageDirection, msg Message) error {
	return nil
}

func (p *localTopicTree) releaseMessage(clientid string, mid uint16, direction MessageDirection) error {
	return nil
}

func (p *localTopicTree) updateMessage(clientid string, mid uint16, direction MessageDirection, state MessageState) {

}
func (p *localTopicTree) addTopic(clientId, topic string, data []byte) {

}

// localTopicTreeFactory
type localTopicTreeFactory struct{}

func newLocalTopicTree(c core.Config) (topicTree, error) {
	d := &localTopicTree{
		config:   c,
		sessions: make(map[string]*Session),
		root: subNode{
			level:    "root",
			children: make(map[string]*subNode),
			subs:     make(map[string]*subLeaf),
		},
	}
	return d, nil
}
