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

type SubscriptionAttribute struct {
	SubscriptionName string    `json:"subscription_name" bson:"SubscriptionName"`
	Subscriber       string    `json:"subscriber" bson:"Subscriber"`
	TopicOwner       string    `json:"topic_owner" bson:"TopicOwner"`
	TopicName        string    `json:"topic_name" bson:"TopicName"`
	Endpoint         string    `json:"endpoint" bson:"Endpoint"`
	NotifyStrategy   string    `json:"notify_strategy" bson:"NotifyStrategy"`
	FilterTag        string    `json:"filter_tag" bson:"FilterTag"`
	CreatedAt        time.Time `json:"created_at" bson:"CreatedAt"`
	LastModifiedAt   time.Time `json:"last_modified_at" bson:"LastModifiedAt"`
}

type Subscription interface {
	Name() string
	GetAttribute() SubscriptionAttribute
	SetAttribute(SubscriptionAttribute)
}
