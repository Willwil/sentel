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
	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/broker/queue"
	"github.com/cloustone/sentel/pkg/config"
)

type subscription struct {
	clientID string
	topic    string
	qos      uint8
	queue    queue.Queue
	retain   bool
}

type topicTree interface {
	// addSubscription add a subscription and context in topic tree
	addSubscription(sub *subscription) error

	// getSubscriptions return client's all subscription topics
	getSubscriptionTopics(clientID string) []string

	// getSubscription return client's subscription by topic
	getSubscription(clientID, topic string) (*subscription, error)

	// removeSubscription remove a subscription from topic tree
	removeSubscription(clientID, topic string) error

	// retainSubscription make a retain process on specified topic
	retainSubscription(clientID, topic string) error

	// addMessage is called when message is published on topic, find the subscriber's queue
	// and write the data to queue that will case queue observer to read data from queue
	addMessage(clientID string, msg *base.Message)

	// retainMessage retain message on specified topic
	retainMessage(clientID string, msg *base.Message)

	// DeleteMessageWithValidator delete message in subdata with condition
	deleteMessageWithValidator(clientID string, validator func(*base.Message) bool)
}

func newTopicTree(c config.Config) (topicTree, error) {
	return newSimpleTopicTree(c)
}
