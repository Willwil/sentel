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

type TopicAttribute struct {
	AccountId              string    `json:"-" bson:"AccountId"`
	TopicName              string    `json:"topic_name" bson:"TopicName"`
	CreatedAt              time.Time `json:"created_at" bson:"CreatedAt"`
	LastModifiedAt         time.Time `json:"last_modified_at" bson:"LastModifiedAt"`
	MaximumMessageSize     uint32    `json:"maximum_message_size" bson:"MaximumMessageSize"`
	MessageRetentionPeriod uint32    `json:"message_retention_period" bson:"MessageRetentionPeriod"`
	LogginEnabled          bool      `json:"loggin_enabled" bson:"LogginEnabled"`
}

type Topic interface {
	Name() string
	SetAttribute(TopicAttribute)
	GetAttribute() TopicAttribute
	Subscribe(filterTag string, notifyStategy string, notifyFmt string) (Subscription, error)
	Unsubscribe(subscriptionId string) error
	SetSubscriptionAttribute(topicName string, subscriptionId string, attr SubscriptionAttribute) error
	GetSubscriptionAttribute(topicName string, subscriptionId string) (SubscriptionAttribute, error)
}

func NewTopic(c config.Config, attr TopicAttribute) (Topic, error) {
	return nil, nil
}
