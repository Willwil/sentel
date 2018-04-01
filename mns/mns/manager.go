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

package mns

import "github.com/cloustone/sentel/pkg/config"

type MnsManager interface {
	// Queue API
	CreateQueue(accountId string, queueName string) (QueueAttribute, error)
	GetQueueAttribute(accountId string, queueName string) (QueueAttribute, error)
	SetQueueAttribute(accountId string, queueName string, attr QueueAttribute) error
	DeleteQueue(accountId string, queueName string) error
	GetQueues(accountId string) []string
	SendQueueMessage(accountId, queueName string, msg QueueMessage) error
	BatchSendQueueMessage(accountId, queueName string, msgs []QueueMessage) error
	ReceiveQueueMessage(accountId, queueName string, waitSeconds int) (QueueMessage, error)
	BatchReceiveQueueMessages(accountId, queueName string, wailtSeconds int, numOfMessages int) ([]QueueMessage, error)
	DeleteQueueMessage(accountId, queueName string, handle string) error
	BatchDeleteQueueMessages(accountId, queueName string, handles []string) error
	PeekQueueMessage(accountId, queueName string, waitSeconds int) (QueueMessage, error)
	BatchPeekQueueMessages(accountId, queueName string, wailtSeconds int, numOfMessages int) ([]QueueMessage, error)
	SetQueueMessageVisibility(accountId, queueName string, handle string, seconds int) error

	// Topic API
	CreateTopic(accountId string, topicName string) (Topic, error)
	GetTopic(accountId string, topicName string) (Topic, error)
	SetTopic(accountId string, topicName string, topic Topic) error
	DeleteTopic(accountId string, topicName string) error
	ListTopics(accountId string) []string

	// Subscription API
	SetSubscription(accountId, topicName string, subscriptionName string, subscription Subscription) error
	GetSubscription(accountId, topicName string, subscriptionName string) (Subscription, error)
	Subscribe(accountId, topicName, subscriptionName, endpoint, filterTag, notifyStrategy, notifiyContentFormat string) error
	Unsubscribe(accountId, topicName string, subscriptionName string) error
	ListTopicSubscriptions(accountId, topicName string, pageno int, pageSize int) ([]Subscription, error)
	PublishMessage(accountId, topicName string, body []byte, tag string, attributes map[string]interface{}) error
}

func NewManager(c config.Config) (MnsManager, error) {
	adaptor, err := NewAdaptor(c)
	return &manager{
		config:  c,
		adaptor: adaptor,
	}, err
}
