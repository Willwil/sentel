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
	"time"

	"github.com/cloustone/sentel/pkg/config"
)

type manager struct {
	config  config.Config
	adaptor Adaptor
}

func (m *manager) CreateQueue(accountId string, queueName string) (QueueAttribute, error) {
	attr := QueueAttribute{
		QueueName:      queueName,
		CreateTime:     time.Now(),
		LastModifyTime: time.Now(),
	}
	err := m.adaptor.AddQueue(accountId, queueName, attr)
	return attr, err
}

func (m *manager) GetQueue(accountId string, queueName string) (Queue, error) {
	if attr, err := m.adaptor.GetQueue(accountId, queueName); err == nil {
		return NewQueue(m.config, attr, m.adaptor)
	}
	return nil, ErrQueueNotExist
}

func (m *manager) DeleteQueue(accountId string, queueName string) error {
	return m.adaptor.RemoveQueue(accountId, queueName)
}

func (m *manager) GetQueues(accountId string) []string {
	return m.adaptor.GetAccountQueues(accountId)
}

func (m *manager) SetQueueAttribute(accountId, queueName string, attr QueueAttribute) error {
	return m.adaptor.UpdateQueue(accountId, queueName, attr)
}

func (m *manager) GetQueueAttribute(accountId, queueName string) (QueueAttribute, error) {
	return m.adaptor.GetQueue(accountId, queueName)
}

func (m *manager) SendQueueMessage(accountId, queueName string, msg QueueMessage) (err error) {
	if queue, err := m.GetQueue(accountId, queueName); err == nil {
		return queue.SendMessage(msg)
	}
	return
}

func (m *manager) BatchSendQueueMessage(accountId, queueName string, msgs []QueueMessage) (err error) {
	if queue, err := m.GetQueue(accountId, queueName); err == nil {
		return queue.BatchSendMessage(msgs)
	}
	return
}

func (m *manager) ReceiveQueueMessage(accountId, queueName string, ws int) (msg QueueMessage, err error) {
	if queue, err := m.GetQueue(accountId, queueName); err == nil {
		return queue.ReceiveMessage(ws)
	}
	return QueueMessage{}, err
}

func (m *manager) BatchReceiveQueueMessages(accountId, queueName string, ws int, numOfMessages int) (msgs []QueueMessage, err error) {
	if queue, err := m.GetQueue(accountId, queueName); err == nil {
		return queue.BatchReceiveMessages(ws, numOfMessages)
	}
	return []QueueMessage{}, err
}

func (m *manager) DeleteQueueMessage(accountId, queueName string, handle string) (err error) {
	if queue, err := m.GetQueue(accountId, queueName); err == nil {
		return queue.DeleteMessage(handle)
	}
	return
}

func (m *manager) BatchDeleteQueueMessages(accountId, queueName string, handles []string) (err error) {
	if queue, err := m.GetQueue(accountId, queueName); err == nil {
		return queue.BatchDeleteMessages(handles)
	}
	return
}

func (m *manager) PeekQueueMessage(accountId, queueName string, ws int) (msg QueueMessage, err error) {
	if queue, err := m.GetQueue(accountId, queueName); err == nil {
		return queue.PeekMessage(ws)
	}
	return QueueMessage{}, err
}

func (m *manager) BatchPeekQueueMessages(accountId, queueName string, ws int, numOfMessages int) (msgs []QueueMessage, err error) {
	if queue, err := m.GetQueue(accountId, queueName); err == nil {
		return queue.BatchPeekMessages(ws, numOfMessages)
	}
	return []QueueMessage{}, err
}

func (m *manager) SetQueueMessageVisibility(accountId, queueName string, handle string, seconds int) (err error) {
	if queue, err := m.GetQueue(accountId, queueName); err == nil {
		return queue.SetMessageVisibility(handle, seconds)
	}
	return
}

// Topic API
func (m *manager) CreateTopic(accountId string, topicName string) (TopicAttribute, error) {
	topicAttr := TopicAttribute{
		TopicName:          topicName,
		CreatedAt:          time.Now(),
		LastModifiedAt:     time.Now(),
		LogginEnabled:      false,
		MaximumMessageSize: 256,
	}
	err := m.adaptor.AddTopic(accountId, topicName, topicAttr)
	return topicAttr, err
}

func (m *manager) GetTopic(accountId string, topicName string) (topic Topic, err error) {
	if topicAttr, err := m.adaptor.GetTopic(accountId, topicName); err == nil {
		return NewTopic(m.config, topicAttr)
	}
	return nil, err
}

func (m *manager) DeleteTopic(accountId string, topicName string) error {
	return m.adaptor.RemoveTopic(accountId, topicName)
}

func (m *manager) ListTopics(account string) []string {
	return m.adaptor.GetAccountTopics(account)
}

func (m *manager) SetTopicAttribute(accountId, topicName string, attr TopicAttribute) error {
	return m.adaptor.UpdateTopic(accountId, topicName, attr)
}

func (m *manager) GetTopicAttribute(accountId, topicName string) (TopicAttribute, error) {
	return m.adaptor.GetTopic(accountId, topicName)
}

// Subscription API
func (m *manager) GetSubscription(accountId, topicName, subscriptionName string) (subscription Subscription, err error) {
	if attr, err := m.adaptor.GetSubscription(accountId, topicName, subscriptionName); err == nil {
		return NewSubscription(m.config, attr)
	}
	return nil, err
}

func (m *manager) SetSubscriptionAttribute(accountId, topicName, subscriptionName string, attr SubscriptionAttribute) error {
	return m.adaptor.UpdateSubscription(accountId, topicName, subscriptionName, attr)
}

func (m *manager) GetSubscriptionAttribute(accountId, topicName string, subscriptionName string) (SubscriptionAttribute, error) {
	return m.adaptor.GetSubscription(accountId, topicName, subscriptionName)
}

func (m *manager) Subscribe(accountId, topicName, subscriptionName, endpoint, filterTag, notifyStrategy, notifyContentFormat string) error {
	attr := SubscriptionAttribute{
		TopicName:           topicName,
		SubscriptionName:    subscriptionName,
		Endpoint:            endpoint,
		FilterTag:           filterTag,
		NotifyStrategy:      notifyStrategy,
		NotifyContentFormat: notifyContentFormat,
		CreatedAt:           time.Now(),
		LastModifiedAt:      time.Now(),
	}
	return m.adaptor.AddSubscription(accountId, topicName, attr)
}

func (m *manager) Unsubscribe(accountId, topicName string, subscriptionName string) error {
	return m.adaptor.RemoveSubscription(accountId, topicName, subscriptionName)
}

func (m *manager) ListTopicSubscriptions(accountId, topicName string, pages int, pageSize int, startIndex int) ([]SubscriptionAttribute, error) {
	subscriptionNames, err := m.adaptor.GetAccountSubscriptionsWithPage(accountId, topicName, pages, pageSize, startIndex)
	subscriptions := []SubscriptionAttribute{}
	for _, subscriptionName := range subscriptionNames {
		attr, _ := m.adaptor.GetSubscription(accountId, topicName, subscriptionName)
		subscriptions = append(subscriptions, attr)
	}
	return subscriptions, err
}

// PublishMessage publish message to all subscribers on the topic
func (m *manager) PublishMessage(accountId, topicName string, body []byte, tag string, attributes map[string]interface{}) error {
	subscriptionNames, err := m.adaptor.GetAccountSubscriptions(accountId, topicName)
	for _, subscriptionName := range subscriptionNames {
		attr, _ := m.adaptor.GetSubscription(accountId, topicName, subscriptionName)
		if endpoint, err := NewEndpoint(m.config, attr); err == nil {
			endpoint.PushMessage(body, tag, attributes)
		}
	}
	return err
}
