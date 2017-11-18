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

package mqtt

import (
	"errors"
	"strings"

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
	retainMsg *StorageMessage
}

type localStorage struct {
	config   core.Config
	sessions map[string]*mqttSession
	root     subNode
}

// Open local storage
func (p *localStorage) Open() error {
	glog.Info("local storage Open")

	return nil
}

// Close local storage
func (p *localStorage) Close() {

}

// Backup serialize local storage
func (p *localStorage) Backup(shutdown bool) error {
	return nil
}

// Restore recover data from serialization
func (p *localStorage) Restore() error {
	return nil
}

// FindSession find session by id
func (p *localStorage) FindSession(id string) (*mqttSession, error) {
	if _, ok := p.sessions[id]; !ok {
		return nil, errors.New("Session id does not exist")
	}

	return p.sessions[id], nil
}

// DeleteSession delete session by id
func (p *localStorage) DeleteSession(id string) error {
	if _, ok := p.sessions[id]; !ok {
		return errors.New("Session id does not exist")
	}
	delete(p.sessions, id)
	return nil
}

// UpdateSession update session
func (p *localStorage) UpdateSession(s *mqttSession) error {
	if _, ok := p.sessions[s.id]; !ok {
		return errors.New("Session id does not exist")
	}
	p.sessions[s.id] = s
	return nil
}

// RegisterSession register new session
func (p *localStorage) RegisterSession(s *mqttSession) error {
	if _, ok := p.sessions[s.id]; ok {
		return errors.New("Session id already exists")
	}

	glog.Infof("RegisterSession: id is %s", s.id)
	p.sessions[s.id] = s
	return nil
}

// Device
// AddDevice
// func (p *localStorage) AddDevice(d Device) error {
// 	return nil
// }

// func (p *localStorage) DeleteDevice(id string) error {
// 	return nil
// }

// func (p *localStorage) UpdateDevice(d Device) error {
// 	return nil
// }

// func (p *localStorage) GetDeviceState(id string) (int, error) {
// 	return 0, nil
// }

// func (p *localStorage) SetDeviceState(state int) error {
// 	return nil
// }

// // Topic
// func (p *localStorage) TopicExist(t Topic) (bool, error) {
// 	return false, nil
// }

// func (p *localStorage) AddTopic(t Topic) error {
// 	return nil
// }

// func (p *localStorage) DeleteTopic(id string) error {
// 	return nil
// }

// func (p *localStorage) UpdateTopic(t Topic) error {
// 	return nil
// }

// func (p *localStorage) AddSubscriber(t Topic, clientid string) error {
// 	return nil
// }

// func (p *localStorage) RemoveSubscriber(t Topic, clientid string) error {
// 	return nil
// }

// func (p *localStorage) GetTopicSubscribers(t Topic) ([]string, error) {
// 	return nil, nil
// }

func (p *localStorage) findNode(node *subNode, lev string) *subNode {
	for k, v := range node.children {
		if k == lev {
			return v
		}
	}

	return nil
}

func (p *localStorage) addNode(node *subNode, lev string) *subNode {
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
func (p *localStorage) AddSubscription(sessionid string, topic string, qos uint8) error {
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
func (p *localStorage) RetainSubscription(sessionid string, topic string, qos uint8) error {
	return nil
}

func (p *localStorage) RemoveSubscription(sessionid string, topic string) error {
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
func (p *localStorage) FindMessage(clientid string, mid uint16) (bool, error) {
	return false, nil
}

func (p *localStorage) StoreMessage(clientid string, msg StorageMessage) error {
	return nil
}

func (p *localStorage) DeleteMessageWithValidator(clientid string, validator func(msg StorageMessage) bool) {

}

func (p *localStorage) DeleteMessage(clientid string, mid uint16, direction MessageDirection) error {
	return nil
}

func (p *localStorage) subProcess(clientid string, msg *StorageMessage, node *subNode, setRetain bool) error {
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

		s.sendPublish(v.qos, msg.Qos, msg.Topic, msg.Payload)
	}
	return nil
}

func (p *localStorage) subSearch(clientid string, msg *StorageMessage, node *subNode, levels []string, setRetain bool) error {
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

func (p *localStorage) QueueMessage(clientid string, msg StorageMessage) error {
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

func (p *localStorage) GetMessageTotalCount(clientid string) int {
	return 0
}

func (p *localStorage) InsertMessage(clientid string, mid uint16, direction MessageDirection, msg StorageMessage) error {
	return nil
}

func (p *localStorage) ReleaseMessage(clientid string, mid uint16, direction MessageDirection) error {
	return nil
}

func (p *localStorage) UpdateMessage(clientid string, mid uint16, direction MessageDirection, state MessageState) {

}

// localStorageFactory
type localStorageFactory struct{}

func (p *localStorageFactory) New(c core.Config) (Storage, error) {
	d := &localStorage{
		config:   c,
		sessions: make(map[string]*mqttSession),
		root: subNode{
			level:    "root",
			children: make(map[string]*subNode),
			subs:     make(map[string]*subLeaf),
		},
	}
	return d, nil
}
