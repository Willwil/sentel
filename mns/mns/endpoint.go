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
	"strings"
	"time"

	"github.com/cloustone/sentel/pkg/config"
)

type Endpoint interface {
	GetAttribute() EndpointAttribute
	PushMessage(body []byte, tag string, attrs map[string]interface{}) error
}

type EndpointAttribute struct {
	Name           string    `json:"endpoint_name" bson:"EndpointName"`
	Type           string    `json:"endpoint_type" bson:"EndpointType"`
	URI            string    `json:"endpoint_uri" bson:"EndpointURI"`
	CreatedAt      time.Time `json:"created_at,omitempty" bson:"CreatedAt,omitempty"`
	LastModifiedAt time.Time `json:"last_modified_at,omitempty" bson:"LastModifiedAt,omitempty"`
}

func NewEndpoint(c config.Config, subscription Subscription) (Endpoint, error) {
	if names := strings.Split(subscription.Endpoint, ":"); len(names) > 0 {
		switch names[0] {
		case "http":
			return newHttpEndpoint(c, subscription)
		case "mns":
			return newQueueEndpoint(c, subscription)
		case "mail":
			return newMailEndpoint(c, subscription)
		case "sms":
			return newSMSEndpoint(c, subscription)
		}
	}
	return nil, ErrInvalidParameter
}
