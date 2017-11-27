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

package subtree

import "github.com/cloustone/sentel/core"

type TopicTree interface {
	Backup(shutdown bool) error
	Restore() error

	// Session
	FindSession(id string) (*Session, error)
	DeleteSession(id string) error
	UpdateSession(s *Session) error
	RegisterSession(s *Session) error

	// Subscription
	AddSubscription(sessionid string, topic string, qos uint8) error
	RetainSubscription(sessionid string, topic string, qos uint8) error
	RemoveSubscription(sessionid string, topic string) error

	// Message Management
	FindMessage(clientid string, mid uint16) (*Message, error)
	StoreMessage(clientid string, msg Message) error
	DeleteMessageWithValidator(clientid string, validator func(Message) bool)
	DeleteMessage(clientid string, mid uint16, direction MessageDirection) error

	QueueMessage(clientid string, msg Message) error
	GetMessageTotalCount(clientid string) int
	InsertMessage(clientid string, mid uint16, direction MessageDirection, msg Message) error
	ReleaseMessage(clientid string, mid uint16, direction MessageDirection) error
	UpdateMessage(clientid string, mid uint16, direction MessageDirection, state MessageState)
}

func newTopicTree(c core.Config) (TopicTree, error) {
	return newLocalTopicTree(c)
}
