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

// Queue API
func CreateQueue(accountId string, queueName string) (Queue, error) {
	return getManager().CreateQueue(accountId, queueName)
}
func GetQueue(accountId string, queueName string) (Queue, error) {
	return getManager().GetQueue(accountId, queueName)
}
func DeleteQueue(accountId string, queueName string) error {
	return getManager().DeleteQueue(accountId, queueName)
}

func GetQueueList(accountId string) ([]string, error) {
	return getManager().GetQueueList(accountId)
}

// Topic API
func CreateTopic(accountId string, topicName string) (Topic, error) {
	return getManager().CreateTopic(accountId, topicName)
}

func GetTopic(accountId string, topicName string) (Topic, error) {
	return getManager().GetTopic(accountId, topicName)
}

func DeleteTopic(accountId string, topicName string) error {
	return getManager().DeleteTopic(accountId, topicName)
}

func ListTopics(accountId string) []string {
	return getManager().ListTopics(accountId)
}

// Subscription API
func GetSubscription(accountId string, subscriptionId string) (Subscription, error) {
	return getManager().GetSubscription(accountId, subscriptionId)
}

func Subscribe(endpoint string, filterTag string, notifyStrategy string, notifyContentFormat string) error {
	return getManager().Subscribe(endpoint, filterTag, notifyStrategy, notifyContentFormat)
}

func Unsubscribe(topicName string, subscriptionId string) error {
	return getManager().Unsubscribe(topicName, subscriptionId)
}

func ListTopicSubscriptions(topicName string, pages uint32, pageSize uint32, startIndex uint32) ([]SubscriptionAttr, error) {
	return getManager().ListTopicSubscriptions(topicName, pages, pageSize, startIndex)
}

func PublishMessage(topicName string, body []byte, tag string, attributes map[string]string) error {
	return getManager().PublishMessage(topicName, body, tag, attributes)
}
