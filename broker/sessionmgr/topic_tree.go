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
	// addSubscription add a subscription and context in topic tree
	addSubscription(clientId, topic string, qos uint8, q queue.Queue) error

	// getSubscription return client's subscription
	getSubscription(clientId string) []string

	// removeSubscription remove a subscription from topic tree
	removeSubscription(clientId, topic string) error

	// retainSubscription make a retain process on specified topic
	retainSubscription(clientId, topic string) error

	// addMessage is called when message is published on topic, find the subscriber's queue
	// and write the data to queue that will case queue observer to read data from queue
	addMessage(clientId, topic string, data []byte)

	// retainMessage retain message on specified topic
	retainMessage(clientId string, msg *Message)

	// DeleteMessageWithValidator delete message in subdata with condition
	deleteMessageWithValidator(clientId string, validator func(*Message) bool)
}

func newTopicTree(c core.Config) (topicTree, error) {
	return newSimpleTopicTree(c)
}
