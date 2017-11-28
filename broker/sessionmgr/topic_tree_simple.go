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

type simpleTopicTree struct {
	root subNode
}

func (p *simpleTopicTree) findNode(node *subNode, lev string) *subNode {
	for k, v := range node.children {
		if k == lev {
			return v
		}
	}
	return nil
}

func (p *simpleTopicTree) addNode(node *subNode, lev string) *subNode {
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
func (p *simpleTopicTree) addSubscription(sessionid string, topic string, qos uint8, q queue.Queue) error {
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
func (p *simpleTopicTree) retainSubscription(sessionid string, topic string) error {
	return nil
}

func (p *simpleTopicTree) removeSubscription(sessionid string, topic string) error {
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

func (p *simpleTopicTree) subProcess(clientid string, msg *Message, node *subNode, setRetain bool) error {
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

func (p *simpleTopicTree) subSearch(clientid string, msg *Message, node *subNode, levels []string, setRetain bool) error {
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

func (p *simpleTopicTree) queueMessage(clientid string, msg Message) error {
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

func (p *simpleTopicTree) addMessage(clientId, topic string, data []byte) {

}

// retainMessage retain message on specified topic
func (p *simpleTopicTree) retainMessage(clientId, msg *Message) {
}

// DeleteMessageWithValidator delete message in subdata with condition
func (p *simpleTopicTree) deleteMessageWithValidator(clientId string, validator func(Message) bool) {
}

// DeleteMessge delete message specified by idfrom subdata
func (p *simpleTopicTree) deleteMessage(clientId string, mid uint16) error {
	return nil
}

// simpleTopicTreeFactory
type simpleTopicTreeFactory struct{}

func newSimpleTopicTree(c core.Config) (topicTree, error) {
	d := &simpleTopicTree{
		root: subNode{
			level:    "root",
			children: make(map[string]*subNode),
			subs:     make(map[string]*subLeaf),
		},
	}
	return d, nil
}
