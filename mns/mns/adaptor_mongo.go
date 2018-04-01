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

func (m *mongoAdaptor) GetQueue(accountId string, queueName string) (QueueAttribute, error) {
	c := m.session.C(collectionQueues)
	attr := QueueAttribute{}
	err := c.Find(bson.M{"AccountId": accountId, "QueueName": queueName}).One(&attr)
	return attr, err
}

func (m *mongoAdaptor) AddQueue(accountId string, queueName string, attr QueueAttribute) error {
	c := m.session.C(collectionQueues)
	attr.AccountId = accountId
	attr.QueueName = queueName
	return c.Insert(attr)
}

func (m *mongoAdaptor) UpdateQueue(accountId string, queueName string, attr QueueAttribute) error {
	c := m.session.C(collectionQueues)
	attr.AccountId = accountId
	attr.QueueName = queueName
	return c.Update(bson.M{"AccountId": accountId, "QueueName": queueName}, attr)
}

func (m *mongoAdaptor) RemoveQueue(accountId string, queueName string) error {
	c := m.session.C(collectionQueues)
	return c.Remove(bson.M{"AccountId": accountId, "queueName": queueName})
}

func (m *mongoAdaptor) GetAccountQueues(accountId string) []string {
	attrs := []QueueAttribute{}
	queues := []string{}
	c := m.session.C(collectionQueues)
	if err := c.Find(bson.M{"AccountId": accountId}).All(attrs); err == nil {
		for _, attr := range attrs {
			queues = append(queues, attr.QueueName)
		}
	}
	return queues
}

func (m *mongoAdaptor) RemoveAccountQueues(accountId string) {
	c := m.session.C(collectionQueues)
	c.Remove(bson.M{"AccountId": accountId})
}

// Topic Object
func (m *mongoAdaptor) GetTopic(accountId string, topicName string) (TopicAttribute, error) {
	c := m.session.C(collectionTopics)
	attr := TopicAttribute{}
	err := c.Find(bson.M{"AccountId": accountId, "TopicName": topicName}).One(&attr)
	return attr, err
}

func (m *mongoAdaptor) AddTopic(accountId string, topicName string, attr TopicAttribute) error {
	c := m.session.C(collectionTopics)
	attr.AccountId = accountId
	attr.TopicName = topicName
	return c.Insert(attr)
}

func (m *mongoAdaptor) UpdateTopic(accountId string, topicName string, attr TopicAttribute) error {
	c := m.session.C(collectionTopics)
	attr.AccountId = accountId
	attr.TopicName = topicName
	return c.Update(bson.M{"AccountId": accountId, "TopicName": topicName}, attr)
}

func (m *mongoAdaptor) RemoveTopic(accountId string, topicName string) error {
	c := m.session.C(collectionTopics)
	return c.Remove(bson.M{"AccountId": accountId, "TopicName": topicName})
}

func (m *mongoAdaptor) GetAccountTopics(accountId string) []string {
	c := m.session.C(collectionTopics)
	attrs := []TopicAttribute{}
	topics := []string{}
	if err := c.Find(bson.M{"AccountId": accountId}).All(&attrs); err == nil {
		for _, attr := range attrs {
			topics = append(topics, attr.TopicName)
		}
	}
	return topics
}

func (m *mongoAdaptor) RemoveAccountTopics(accountId string) {
	c := m.session.C(collectionTopics)
	c.Remove(bson.M{"AccountId": accountId})
}

// Subscription Object
func (m *mongoAdaptor) GetSubscription(accountId string, topicName string, subscriptionName string) (SubscriptionAttribute, error) {
	c := m.session.C(collectionSubscriptions)
	attr := SubscriptionAttribute{}
	err := c.Find(bson.M{"AccountId": accountId, "TopicName": topicName, "SubscriptionName": subscriptionName}).One(&attr)
	return attr, err
}

func (m *mongoAdaptor) AddSubscription(accountId string, topicName string, attr SubscriptionAttribute) error {
	c := m.session.C(collectionSubscriptions)
	attr.TopicName = topicName
	attr.AccountId = accountId
	return c.Insert(attr)
}

func (m *mongoAdaptor) UpdateSubscription(accountId string, topicName string, subscriptionName string, attr SubscriptionAttribute) error {
	c := m.session.C(collectionSubscriptions)
	attr.TopicName = topicName
	attr.AccountId = accountId
	attr.SubscriptionName = subscriptionName
	return c.Update(bson.M{"AccountId": accountId, "TopicName": topicName, "SubscriptionName": subscriptionName}, attr)
}

func (m *mongoAdaptor) RemoveSubscription(accountId string, topicName string, subscriptionName string) error {
	c := m.session.C(collectionSubscriptions)
	return c.Remove(bson.M{"AccountId": accountId, "TopicName": topicName, "subscriptionName": subscriptionName})
}

func (m *mongoAdaptor) GetAccountSubscriptions(accountId string, topicName string) ([]string, error) {
	c := m.session.C(collectionSubscriptions)
	attrs := []SubscriptionAttribute{}
	subscriptions := []string{}
	err := c.Find(bson.M{"AccountId": accountId, "TopicName": topicName}).All(&attrs)
	for _, attr := range attrs {
		subscriptions = append(subscriptions, attr.SubscriptionName)
	}
	return subscriptions, err
}

func (m *mongoAdaptor) GetAccountSubscriptionsWithPage(accountId string, topicName string, pages int, pageSize int, startIndex int) ([]string, error) {
	return []string{}, errors.New("Not implemented")
}

func (m *mongoAdaptor) RemoveAccountSubscriptions(accountId string, topicName string) {
	c := m.session.C(collectionSubscriptions)
	c.Remove(bson.M{"AccountId": accountId, "TopicName": topicName})
}
