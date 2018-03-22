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

type MnsManager interface {
	// Queue API
	CreateQueue(accountId string, queueName string) (QueueAttribute, error)
	GetQueue(accountId string, queueName string) (Queue, error)
	DeleteQueue(accountId string, queueName string) error
	GetQueues(accountId string) ([]string, error)

	// Topic API
	CreateTopic(accountId string, topicName string) (Topic, error)
	GetTopic(accountId string, topicName string) (Topic, error)
	DeleteTopic(accountId string, topicName string) error
	ListTopics(accountId string) []string

	// Subscription API
	GetSubscription(accountId string, subscriptionId string) (Subscription, error)
	Subscribe(accountId, subscriptionName, endpoint, filterTag, notifyStrategy, notifiyContentFormat string) error
	Unsubscribe(accountId, topicName string, subscriptionId string) error
	ListTopicSubscriptions(accountId, topicName string, pages int, pageSize int, startIndex int) ([]SubscriptionAttribute, error)
	PublishMessage(accountId, topicName string, body []byte, tag string, attributes map[string]string) error
}

func NewManager(c config.Config) (MnsManager, error) {
	adaptor, err := NewAdaptor(c)
	return &manager{
		config:  c,
		adaptor: adaptor,
	}, err
}

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
	return nil, ErrInternalError
}

func (m *manager) DeleteQueue(accountId string, queueName string) error { return ErrInternalError }

func (m *manager) GetQueues(accountId string) ([]string, error) {
	return []string{}, nil
}

// Topic API
func (m *manager) CreateTopic(accountId string, topicName string) (Topic, error) {
	return nil, ErrInternalError
}

func (m *manager) GetTopic(accountId string, topicName string) (Topic, error) {
	return nil, ErrInternalError
}

func (m *manager) DeleteTopic(accountId string, topicName string) error { return nil }
func (m *manager) ListTopics(account string) []string                   { return []string{} }

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
