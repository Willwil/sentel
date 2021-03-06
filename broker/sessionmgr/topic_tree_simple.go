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
	"fmt"
	"strings"
	"sync"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/broker/queue"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/golang/glog"
)

// usbNode is subscription node in topic tree
type topicNode struct {
	level    string                   // Topic level for the node
	children map[string]*topicNode    // Childern topic node
	subs     map[string]*subscription // All subscription contexts
	msgs     []*base.Message          // All retained messages
}

// simpleTopicTree manage all subscripted topic
type simpleTopicTree struct {
	root   topicNode           // Root node
	mutex  sync.Mutex          // Mutex for concurrence context
	topics map[string][]string // All clients and subsribned topics
}

// findNode return sub node with specified topic
func (p *simpleTopicTree) findNode(node *topicNode, topic string) *topicNode {
	if n, found := node.children[topic]; found {
		return n
	}
	return nil
}

// addNode add a sub node in root node
func (p *simpleTopicTree) addNode(node *topicNode, level string, q queue.Queue) *topicNode {
	if _, found := node.children[level]; !found {
		n := &topicNode{
			level:    level,
			children: make(map[string]*topicNode),
			subs:     make(map[string]*subscription),
			msgs:     []*base.Message{},
		}
		node.children[level] = n
	}
	return node.children[level]
}

// addSubscription add subsription into topic tree
func (p *simpleTopicTree) addSubscription(sub *subscription) error {
	glog.Infof("topictree:addSubscription: clientID is %s, topic is %s, qos is %d",
		sub.clientID, sub.topic, sub.qos)

	// Get topic slice
	topics := strings.Split(sub.topic, "/")
	if len(topics) == 0 {
		return errors.New("Invalid topic")
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()

	node := &p.root
	for _, level := range topics {
		node = p.addNode(node, level, sub.queue)
	}

	node.subs[sub.clientID] = sub
	if _, found := p.topics[sub.clientID]; !found {
		p.topics[sub.clientID] = []string{sub.topic}
	} else {
		p.topics[sub.clientID] = append(p.topics[sub.clientID], sub.topic)
	}
	return nil
}

// retainSubscription retain the subscription
func (p *simpleTopicTree) retainSubscription(clientID string, topic string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	node := &p.root
	topics := strings.Split(topic, "/")
	for _, level := range topics {
		node = p.findNode(node, level)
		if node == nil {
			return fmt.Errorf("topic tree: invalid topic '%s'", topic)
		}
	}
	if sub, found := node.subs[clientID]; found {
		sub.retain = true
		return nil
	}
	return fmt.Errorf("topic tree: invalid client id '%s'", clientID)
}

// removeSubscription remove subscription from topic tree
func (p *simpleTopicTree) removeSubscription(clientID string, topic string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	node := &p.root
	levels := strings.Split(topic, "/")
	for _, level := range levels {
		node = p.findNode(node, level)
		if node == nil {
			return fmt.Errorf("topic tree: invalid subscription '%s'", topic)
		}
	}
	if _, ok := node.subs[clientID]; !ok {
		return fmt.Errorf("topic tree: invalid client id '%s'", clientID)
	}
	delete(node.subs, clientID)

	// Remove the topic
	if _, found := p.topics[clientID]; found {
		list := p.topics[clientID]
		for i, t := range list {
			if t == topic {
				list = append(list[:i], list[i+1:]...)
				break
			}
		}
	}
	return nil
}

// searchNode search specified node recursively
func (p *simpleTopicTree) searchNode(node *topicNode, levels []string, setRetain bool) *topicNode {
	for k, v := range node.children {
		sr := setRetain
		if len(levels) != 0 && (k == levels[0] || k == "+") {
			if k == "+" {
				sr = false
			}
			ss := levels[1:]
			if len(ss) == 0 {
				return v
			}
			return p.searchNode(v, ss, sr)
		} else if k == "#" && len(v.children) > 0 {
			return v
		}
	}
	return nil
}

// addMessage publish a message on topic tree
func (p *simpleTopicTree) addMessage(clientID string, msg *base.Message) {
	levels := strings.Split(msg.Topic, "/")
	p.mutex.Lock()
	defer p.mutex.Unlock()
	node := p.searchNode(&p.root, levels, true)
	if node != nil {
		for _, sub := range node.subs {
			glog.Infof("topic tree: publishing message with client id '%s'", clientID)
			newMsg := &base.Message{
				Topic:     msg.Topic,
				PacketId:  msg.PacketId,
				Direction: msg.Direction,
				State:     msg.State,
				Qos:       sub.qos,
				Dup:       false,
				Retain:    msg.Retain,
				Time:      msg.Time,
				Payload:   msg.Payload,
			}
			sub.queue.Pushback(newMsg)
		}
	}
}

// retainMessage retain message on specified topic
func (p *simpleTopicTree) retainMessage(clientID string, msg *base.Message) {
	levels := strings.Split(msg.Topic, "/")
	p.mutex.Lock()
	defer p.mutex.Unlock()
	node := p.searchNode(&p.root, levels, true)
	if node != nil {
		node.msgs = append(node.msgs, msg)
	} else {
		glog.Errorf("topic tree: Failed to retain message for '%s'", clientID)
	}
}

// deleteMessageWithValidator delete message in subdata with condition
func (p *simpleTopicTree) deleteMessageWithValidator(clientID string, validator func(*base.Message) bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, topic := range p.topics[clientID] {
		levels := strings.Split(topic, "/")
		node := p.searchNode(&p.root, levels, true)
		for index, msg := range node.msgs {
			if validator(msg) {
				node.msgs = append(node.msgs[:index], node.msgs[index+1:]...)
			}
		}
	}
}

// getSubscriptions return client's all subscription topics
func (p *simpleTopicTree) getSubscriptionTopics(clientID string) []string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.topics[clientID]
}

// getSubscription return client's subscription by topic
func (p *simpleTopicTree) getSubscription(clientID, topic string) (*subscription, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	node := &p.root
	levels := strings.Split(topic, "/")
	for _, level := range levels {
		node = p.findNode(node, level)
		if node == nil {
			return nil, fmt.Errorf("topic tree: invalid subscription '%s'", topic)
		}
	}
	if _, ok := node.subs[clientID]; !ok {
		return nil, fmt.Errorf("topic tree: invalid client id '%s'", clientID)
	}
	return node.subs[clientID], nil
}

func newSimpleTopicTree(c config.Config) (topicTree, error) {
	d := &simpleTopicTree{
		root: topicNode{
			level:    "root",
			children: make(map[string]*topicNode),
		},
		topics: make(map[string][]string),
	}
	return d, nil
}
