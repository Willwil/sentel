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
	"github.com/cloustone/sentel/broker/queue"
	"github.com/cloustone/sentel/core"
)

type TopicTree interface {
	// Session
	FindSession(id string) (*Session, error)
	DeleteSession(id string) error
	UpdateSession(s *Session) error
	RegisterSession(s *Session) error

	// Subscription
	AddSubscription(clientId string, topic string, qos uint8, q queue.Queue) error
	RetainSubscription(clientId string, topic string, qos uint8) error
	RemoveSubscription(clientId string, topic string) error

	// Message Management
	FindMessage(clientId string, mid uint16) (*Message, error)
	StoreMessage(clientId string, msg Message) error
	DeleteMessageWithValidator(clientId string, validator func(Message) bool)
	DeleteMessage(clientId string, mid uint16, direction MessageDirection) error

	QueueMessage(clientId string, msg Message) error
	GetMessageTotalCount(clientId string) int
	InsertMessage(clientId string, mid uint16, direction MessageDirection, msg Message) error
	ReleaseMessage(clientId string, mid uint16, direction MessageDirection) error
	UpdateMessage(clientId string, mid uint16, direction MessageDirection, state MessageState)

	// Topic
	// AddTopic is called when topic is published, find the subscriber's queue
	// and write the data to queue that will case queue observer to read data
	// from queue
	AddTopic(clientId, topic string, data []byte)
}

func newTopicTree(c core.Config) (TopicTree, error) {
	return newLocalTopicTree(c)
}
