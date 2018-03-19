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
	"errors"

	"github.com/cloustone/sentel/pkg/config"
)

type MnsManager interface {
	// Queue API
	CreateQueue(accountId string, queueName string) (Queue, error)
	GetQueue(accountId string, queueName string) (Queue, error)
	DeleteQueue(accountId string, queueName string) error
	GetQueueList(accountId string) ([]string, error)

	// Topic API
	CreateTopic(accountId string, topicName string) (Topic, error)
	GetTopic(accountId string, topicName string) (Topic, error)
	DeleteTopic(accountId string, topicName string) error
	ListTopics(account string) []string

	// Subscription API
	GetSubscription(accountId string, subscriptionId string) (Subscription, error)
	Subscribe(endpoint string, filterTag string, notifyStrategy string, notifiyContentFormat string) error
	Unsubscribe(topicName string, subscriptionId string) error
	ListTopicSubscriptions(topicName string, pages uint32, pageSize uint32, startIndex uint32) ([]SubscriptionAttr, error)
	PublishMessage(topicName string, body []byte, tag string, attributes map[string]string) error
}

func NewManager(c config.Config) (MnsManager, error) {
	return nil, ErrNotImplemented
}

type manager struct {
}

func (m *manager) CreateQueue(accountId string, queueName string) (Queue, error) {
	return nil, ErrNotImplemented
}
func (m *manager) GetQueue(accountId string, queueName string) (Queue, error) {
	return nil, ErrNotImplemented
}
func (m *manager) DeleteQueue(accountId string, queueName string) error { return ErrNotImplemented }

// Topic API
func (m *manager) CreateTopic(accountId string, topicName string) (Topic, error) {
	return nil, errors.New("Not implmented")
}
func (m *manager) GetTopic(accountId string, topicName string) (Topic, error) {
	return nil, errors.New("Not implmented")
}
func (m *manager) DeleteTopic(accountId string, topicName string) error { return nil }
func (m *manager) ListTopics(account string) []string                   { return []string{} }

// Subscription API
func (m *manager) GetSubscription(accountId string, subscriptionId string) (Subscription, error) {
	return nil, ErrNotImplemented
}
func (m *manager) Subscribe(endpoint string, filterTag string, notifyStrategy string, notifiyContentFormat string) error {
	return ErrNotImplemented
}

func (m *manager) Unsubscribe(topicName string, subscriptionId string) error {
	return ErrNotImplemented
}

func (m *manager) ListTopicSubscriptions(topicName string, pages uint32, pageSize uint32, startIndex uint32) ([]SubscriptionAttr, error) {
	return []SubscriptionAttr{}, ErrNotImplemented
}

func (m *manager) PublishMessage(topicName string, body []byte, tag string, attributes map[string]string) error {
	return ErrNotImplemented
}
