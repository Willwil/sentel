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
	return nil
}
func (m *manager) GetQueueAttribute(accountId, queueName string) (QueueAttribute, error) {
	return QueueAttribute{}, nil
}
func (m *manager) SendQueueMessage(accountId, queueName string, msg Message) error {
	return nil
}
func (m *manager) BatchSendQueueMessage(accountId, queueName string, msgs []Message) error {
	return nil
}
func (m *manager) ReceiveQueueMessage(accountId, queueName string, waitSeconds int) (Message, error) {
	return Message{}, nil
}
func (m *manager) BatchReceiveQueueMessages(accountId, queueName string, wailtSeconds int, numOfMessages int) ([]Message, error) {
	return nil, nil
}

func (m *manager) DeleteQueueMessage(accountId, queueName string, handle string) error {
	return nil
}
func (m *manager) BatchDeleteQueueMessages(accountId, queueName string, handles []string) error {
	return nil
}
func (m *manager) PeekQueueMessage(accountId, queueName string, waitSeconds int) (Message, error) {
	return Message{}, nil
}
func (m *manager) BatchPeekQueueMessages(accountId, queueName string, wailtSeconds int, numOfMessages int) ([]Message, error) {
	return nil, nil
}
func (m *manager) SetQueueMessageVisibility(accountId, queueName string, handle string, seconds int) {}

// Topic API
func (m *manager) CreateTopic(accountId string, topicName string) (Topic, error) {
	return nil, ErrInternalError
}

func (m *manager) GetTopic(accountId string, topicName string) (Topic, error) {
	return nil, ErrInternalError
}

func (m *manager) DeleteTopic(accountId string, topicName string) error { return nil }
func (m *manager) ListTopics(account string) []string                   { return []string{} }
func (m *manager) SetTopicAttribute(accountId, topicName string, attr TopicAttribute) error {
	return nil
}
func (m *manager) GetTopicAttribute(accountId, topicName string) (TopicAttribute, error) {
	return TopicAttribute{}, nil
}
func (m *manager) SetSubscriptionAttribute(accountId, topicName string, subscriptionId string, attr SubscriptionAttribute) error {
	return nil
}
func (m *manager) GetSubscriptionAttribute(accountId, topicName string, subscriptionId string) (SubscriptionAttribute, error) {
	return SubscriptionAttribute{}, nil
}

// Subscription API
func (m *manager) GetSubscription(accountId string, subscriptionId string) (Subscription, error) {
	return nil, ErrInternalError
}

func (m *manager) Subscribe(accountId, subscriptionName, endpoint, filterTag, notifyStrategy, notifiyContentFormat string) error {
	return ErrInternalError
}

func (m *manager) Unsubscribe(accountId, topicName string, subscriptionName string) error {
	return ErrInternalError
}

func (m *manager) ListTopicSubscriptions(accountId, topicName string, pages int, pageSize int, startIndex int) ([]SubscriptionAttribute, error) {
	return []SubscriptionAttribute{}, ErrInternalError
}

func (m *manager) PublishMessage(accountId, topicName string, body []byte, tag string, attributes map[string]string) error {
	return ErrInternalError
}
