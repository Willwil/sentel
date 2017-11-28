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

type topicTree interface {
	// Session
	findSession(id string) (*Session, error)
	deleteSession(id string) error
	updateSession(s *Session) error
	registerSession(s *Session) error

	// Subscription
	addSubscription(clientId string, topic string, qos uint8, q queue.Queue) error
	retainSubscription(clientId string, topic string, qos uint8) error
	removeSubscription(clientId string, topic string) error

	// Message Management
	findMessage(clientId string, mid uint16) (*Message, error)
	storeMessage(clientId string, msg Message) error
	deleteMessageWithValidator(clientId string, validator func(Message) bool)
	deleteMessage(clientId string, mid uint16, direction MessageDirection) error

	queueMessage(clientId string, msg Message) error
	getMessageTotalCount(clientId string) int
	insertMessage(clientId string, mid uint16, direction MessageDirection, msg Message) error
	releaseMessage(clientId string, mid uint16, direction MessageDirection) error
	updateMessage(clientId string, mid uint16, direction MessageDirection, state MessageState)

	// Topic
	// AddTopic is called when topic is published, find the subscriber's queue
	// and write the data to queue that will case queue observer to read data
	// from queue
	addTopic(clientId, topic string, data []byte)
}

func newTopicTree(c core.Config) (topicTree, error) {
	return newLocalTopicTree(c)
}
