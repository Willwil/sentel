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
	"fmt"
	"time"

	"github.com/cloustone/sentel/pkg/config"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DBNAME                  = "mns"
	collectionQueues        = "queues"
	collectionTopics        = "topics"
	collectionSubscriptions = "subscriptions"
)

type mongoAdaptor struct {
	config  config.Config
	dbconn  *mgo.Session
	session *mgo.Database
}

func newMongoAdaptor(c config.Config) (Adaptor, error) {
	addr := c.MustString("mongo")
	dbc, err := mgo.DialWithTimeout(addr, 1*time.Second)
	if err != nil {
		return nil, fmt.Errorf("mongo adaptor: '%s'", err.Error())
	}
	return &mongoAdaptor{
		config:  c,
		dbconn:  dbc,
		session: dbc.DB(DBNAME),
	}, nil
}

func (m *mongoAdaptor) GetQueue(accountId string, queueName string) (Queue, error) {
	c := m.session.C(collectionQueues)
	queue := Queue{}
	err := c.Find(bson.M{"AccountId": accountId, "QueueName": queueName}).One(&queue)
	return queue, err
}

func (m *mongoAdaptor) AddQueue(accountId string, queueName string, queue Queue) error {
	c := m.session.C(collectionQueues)
	queue.AccountId = accountId
	queue.QueueName = queueName
	return c.Insert(queue)
}

func (m *mongoAdaptor) UpdateQueue(accountId string, queueName string, queue Queue) error {
	c := m.session.C(collectionQueues)
	queue.AccountId = accountId
	queue.QueueName = queueName
	return c.Update(bson.M{"AccountId": accountId, "QueueName": queueName}, queue)
}

func (m *mongoAdaptor) RemoveQueue(accountId string, queueName string) error {
	c := m.session.C(collectionQueues)
	return c.Remove(bson.M{"AccountId": accountId, "queueName": queueName})
}

func (m *mongoAdaptor) GetAccountQueues(accountId string) []string {
	queues := []Queue{}
	queueNames := []string{}
	c := m.session.C(collectionQueues)
	if err := c.Find(bson.M{"AccountId": accountId}).All(queues); err == nil {
		for _, queue := range queues {
			queueNames = append(queueNames, queue.QueueName)
		}
	}
	return queueNames
}

func (m *mongoAdaptor) RemoveAccountQueues(accountId string) {
	c := m.session.C(collectionQueues)
	c.Remove(bson.M{"AccountId": accountId})
}

// Topic Object
func (m *mongoAdaptor) GetTopic(accountId string, topicName string) (Topic, error) {
	c := m.session.C(collectionTopics)
	attr := Topic{}
	err := c.Find(bson.M{"AccountId": accountId, "TopicName": topicName}).One(&attr)
	return attr, err
}

func (m *mongoAdaptor) AddTopic(accountId string, topicName string, topic Topic) error {
	c := m.session.C(collectionTopics)
	topic.AccountId = accountId
	topic.TopicName = topicName
	return c.Insert(topic)
}

func (m *mongoAdaptor) UpdateTopic(accountId string, topicName string, topic Topic) error {
	c := m.session.C(collectionTopics)
	topic.AccountId = accountId
	topic.TopicName = topicName
	return c.Update(bson.M{"AccountId": accountId, "TopicName": topicName}, topic)
}

func (m *mongoAdaptor) RemoveTopic(accountId string, topicName string) error {
	c := m.session.C(collectionTopics)
	return c.Remove(bson.M{"AccountId": accountId, "TopicName": topicName})
}

func (m *mongoAdaptor) GetAccountTopics(accountId string) []string {
	c := m.session.C(collectionTopics)
	topics := []Topic{}
	topicNames := []string{}
	if err := c.Find(bson.M{"AccountId": accountId}).All(&topics); err == nil {
		for _, topic := range topics {
			topicNames = append(topicNames, topic.TopicName)
		}
	}
	return topicNames
}

func (m *mongoAdaptor) RemoveAccountTopics(accountId string) {
	c := m.session.C(collectionTopics)
	c.Remove(bson.M{"AccountId": accountId})
}

// Subscription Object
func (m *mongoAdaptor) GetSubscription(accountId string, topicName string, subscriptionName string) (Subscription, error) {
	c := m.session.C(collectionSubscriptions)
	attr := Subscription{}
	err := c.Find(bson.M{"AccountId": accountId, "TopicName": topicName, "SubscriptionName": subscriptionName}).One(&attr)
	return attr, err
}

func (m *mongoAdaptor) AddSubscription(accountId string, topicName string, subscription Subscription) error {
	c := m.session.C(collectionSubscriptions)
	subscription.TopicName = topicName
	subscription.AccountId = accountId
	return c.Insert(subscription)
}

func (m *mongoAdaptor) UpdateSubscription(accountId string, topicName string, subscriptionName string, subscription Subscription) error {
	c := m.session.C(collectionSubscriptions)
	subscription.TopicName = topicName
	subscription.AccountId = accountId
	subscription.SubscriptionName = subscriptionName
	return c.Update(bson.M{"AccountId": accountId, "TopicName": topicName, "SubscriptionName": subscriptionName}, subscription)
}

func (m *mongoAdaptor) RemoveSubscription(accountId string, topicName string, subscriptionName string) error {
	c := m.session.C(collectionSubscriptions)
	return c.Remove(bson.M{"AccountId": accountId, "TopicName": topicName, "subscriptionName": subscriptionName})
}

func (m *mongoAdaptor) GetAccountSubscriptions(accountId string, topicName string) ([]string, error) {
	c := m.session.C(collectionSubscriptions)
	subscriptions := []Subscription{}
	subscriptionNames := []string{}
	err := c.Find(bson.M{"AccountId": accountId, "TopicName": topicName}).All(&subscriptions)
	for _, subscription := range subscriptions {
		subscriptionNames = append(subscriptionNames, subscription.SubscriptionName)
	}
	return subscriptionNames, err
}

func (m *mongoAdaptor) GetAccountSubscriptionsWithPage(accountId string, topicName string, pageno int, pageSize int) ([]string, error) {
	c := m.session.C(collectionSubscriptions)
	subscriptions := []Subscription{}
	subscriptionNames := []string{}
	err := c.Find(bson.M{"AccountId": accountId, "TopicName": topicName}).Skip(pageno * pageSize).Limit(pageSize).All(&subscriptions)
	for _, subscription := range subscriptions {
		subscriptionNames = append(subscriptionNames, subscription.SubscriptionName)
	}
	return subscriptionNames, err
}

func (m *mongoAdaptor) RemoveAccountSubscriptions(accountId string, topicName string) {
	c := m.session.C(collectionSubscriptions)
	c.Remove(bson.M{"AccountId": accountId, "TopicName": topicName})
}
