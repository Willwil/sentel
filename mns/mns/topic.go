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

import "time"

type TopicAttribute struct {
	TopicName              string
	CreatedAt              time.Time
	LastModifiedAt         time.Time
	MaximumMessageSize     uint32
	MessageRetentionPeriod uint32
	LogginEnabled          bool
}

type Topic interface {
	Name() string
	SetTopicAttribute(TopicAttribute)
	GetTopicAttribute() TopicAttribute
	Subscribe(filterTag string, notifyStategy string, notifyFmt string) (Subscription, error)
	Unsubscribe(subscriptionId string) error
	SetSubscriptionAttribute(topicName string, subscriptionId string, attr SubscriptionAttr) error
	GetSubscriptionAttribute(topicName string, subscriptionId string) (SubscriptionAttr, error)
}
