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

	"github.com/cloustone/sentel/broker/queue"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

// context is subscription context
type context struct {
	qos      uint8
	queue    queue.Queue // Topic subscriber
	retain   bool
	clientId string
}

// usbNode is subscription node in topic tree
type subNode struct {
	level    string              // Topic level for the node
	children map[string]*subNode // Childern subscription topic
	ctxs     map[string]*context // All subscription contexts
	msgs     []*Message          // All retained messages
}

// simpleTopicTree manage all subscripted topic
type simpleTopicTree struct {
	root   subNode             // Root node
	mutex  sync.Mutex          // Mutex for concurrence context
	topics map[string][]string // All clients and subsribned topics
}

// findNode return sub node with specified topic
func (p *simpleTopicTree) findNode(node *subNode, topic string) *subNode {
	if n, found := node.children[topic]; found {
		return n
	}
	return nil
}

// addNode add a sub node in root node
func (p *simpleTopicTree) addNode(node *subNode, level string, q queue.Queue) *subNode {
	if _, found := node.children[level]; !found {
		n := &subNode{
			level:    level,
			children: make(map[string]*subNode),
			ctxs:     make(map[string]*context),
			msgs:     []*Message{},
		}
		node.children[level] = n
	}
	return node.children[level]
}

// addSubscription add subsription into topic tree
func (p *simpleTopicTree) addSubscription(clientId string, topic string, qos uint8, q queue.Queue) error {
	glog.Infof("topictree:addSubscription: clientId is %s, topic is %s, qos is %d", clientId, topic, qos)

	// Get topic slice
	topics := strings.Split(topic, "/")
	if len(topics) == 0 {
		return errors.New("Invalid topic")
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()

	node := &p.root
	for _, level := range topics {
		node = p.addNode(node, level, q)
	}

	node.ctxs[clientId] = &context{
		qos:   qos,
		queue: q,
	}
	if _, found := p.topics[clientId]; !found {
		p.topics[clientId] = []string{topic}
	} else {
		p.topics[clientId] = append(p.topics[clientId], topic)
	}
	return nil
}

// retainSubscription retain the subscription
func (p *simpleTopicTree) retainSubscription(clientId string, topic string) error {
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
	if ctx, found := node.ctxs[clientId]; found {
		ctx.retain = true
		return nil
	}
	return fmt.Errorf("topic tree: invalid client id '%s'", clientId)
}

// removeSubscription remove subscription from topic tree
func (p *simpleTopicTree) removeSubscription(clientId string, topic string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	node := &p.root
	topics := strings.Split(topic, "/")
	for _, level := range topics {
		node = p.findNode(node, level)
		if node == nil {
			return fmt.Errorf("topic tree: invalid subscription '%s'", topic)
		}
	}
	if _, ok := node.ctxs[clientId]; !ok {
		return fmt.Errorf("topic tree: invalid client id '%s'", clientId)
	}
	delete(node.ctxs, clientId)
	// Remove the topic
	if _, found := p.topics[clientId]; found {
		list := p.topics[clientId]
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
func (p *simpleTopicTree) searchNode(node *subNode, levels []string, setRetain bool) *subNode {
	for k, v := range node.children {
		sr := setRetain
		if len(levels) != 0 && (k == levels[0] || k == "+") {
			if k == "+" {
				sr = false
			}
			ss := levels[1:]
			p.searchNode(v, ss, sr)
			if len(ss) == 0 {
				return v
			}
		} else if k == "#" && len(v.children) > 0 {
			return v
		}
	}
	return nil
}

// addMessage publish a message on topic tree
func (p *simpleTopicTree) addMessage(clientId, topic string, data []byte) {
	levels := strings.Split(topic, "/")
	p.mutex.Lock()
	defer p.mutex.Unlock()
	node := p.searchNode(&p.root, levels, true)
	if node != nil {
		for _, ctx := range node.ctxs {
			glog.Infof("topic tree: publishing message with client id '%s'", clientId)
			ctx.queue.Write(data)
		}
	}
}

// retainMessage retain message on specified topic
func (p *simpleTopicTree) retainMessage(clientId string, msg *Message) {
	levels := strings.Split(msg.Topic, "/")
	p.mutex.Lock()
	defer p.mutex.Unlock()
	node := p.searchNode(&p.root, levels, true)
	if node != nil {
		node.msgs = append(node.msgs, msg)
	} else {
		glog.Fatalf("topic tree: Failed to retain message for '%s'", clientId)
	}
}

// deleteMessageWithValidator delete message in subdata with condition
func (p *simpleTopicTree) deleteMessageWithValidator(clientId string, validator func(*Message) bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, topic := range p.topics[clientId] {
		levels := strings.Split(topic, "/")
		node := p.searchNode(&p.root, levels, true)
		for index, msg := range node.msgs {
			if validator(msg) {
				node.msgs = append(node.msgs[:index], node.msgs[index+1:]...)
			}
		}
	}
}

// getSubscription return client's subscription
func (p *simpleTopicTree) getSubscription(clientId string) []string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.topics[clientId]
}

func newSimpleTopicTree(c core.Config) (topicTree, error) {
	d := &simpleTopicTree{
		root: subNode{
			level:    "root",
			children: make(map[string]*subNode),
		},
		topics: make(map[string][]string),
	}
	return d, nil
}
