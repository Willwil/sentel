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

import (
	"github.com/cloustone/sentel/pkg/config"
)

type Adaptor interface {
	// Queue Object
	GetQueue(accountId string, queueName string) (QueueAttribute, error)
	AddQueue(accountId string, queueName string, attr QueueAttribute) error
	UpdateQueue(accountId string, queueName string, attr QueueAttribute) error
	RemoveQueue(accountId string, queueName string) error
	GetAccountQueues(accountId string) []string
	RemoveAccountQueues(accountId string)
	// Topic Object
	GetTopic(accountId string, topicName string) (TopicAttribute, error)
	AddTopic(accountId string, topicName string, attr TopicAttribute) error
	UpdateTopic(accountId string, topicName string, attr TopicAttribute) error
	RemoveTopic(accountId string, topicName string) error
	GetAccountTopics(accountId string) []string
	RemoveAccountTopics(accountId string)
	// Subscription Object
	GetSubscription(accountId string, topicName string, subscriptionName string) (SubscriptionAttribute, error)
	AddSubscription(accountId string, topicName string, attr SubscriptionAttribute) error
	UpdateSubscription(accountId string, topicName string, subscriptionName string, attr SubscriptionAttribute) error
	RemoveSubscription(accountId string, topicName string, subscriptionName string) error
	GetAccountSubscriptions(accountId string, topicName string) ([]string, error)
	GetAccountSubscriptionsWithPage(accountId string, topicName string, pages int, pageSize int, startIndex int) ([]string, error)
	RemoveAccountSubscriptions(accountId string, topicName string)
}

func NewAdaptor(c config.Config) (Adaptor, error) {
	return &mongoAdaptor{config: c}, ErrInternalError
}
