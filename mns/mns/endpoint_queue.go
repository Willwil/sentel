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
	"strings"
	"time"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
)

// queueEndpoint has scheme "mns:AccountId:queues/queueName"
type queueEndpoint struct {
	config    config.Config
	accountId string
	queueName string
	attribute EndpointAttribute
}

func newQueueEndpoint(c config.Config, subscription Subscription) (endpoint Endpoint, err error) {
	accountId, queueName, err := parseQueueScheme(subscription.Endpoint)
	endpoint = &queueEndpoint{
		config:    c,
		accountId: accountId,
		queueName: queueName,
		attribute: EndpointAttribute{
			Name:           subscription.Endpoint,
			Type:           "mns",
			URI:            subscription.Endpoint,
			CreatedAt:      time.Now(),
			LastModifiedAt: time.Now(),
		},
	}
	return
}

func parseQueueScheme(uri string) (string, string, error) {
	accountId := ""
	queueName := ""
	err := ErrInvalidParameter
	if names := strings.Split(uri, ":"); len(names) == 3 {
		if items := strings.Split(names[2], "/"); len(items) == 2 {
			accountId = names[1]
			queueName = items[1]
			err = nil
		}
	}
	return accountId, queueName, err.With(uri)
}

func (q queueEndpoint) GetAttribute() EndpointAttribute { return q.attribute }
func (q queueEndpoint) PushMessage(body []byte, tag string, attrs map[string]interface{}) error {
	msg := &QueueMessage{
		TopicName:   fmt.Sprintf(FormatOfQueueName, q.accountId, q.queueName),
		MessageBody: body,
		MessageTag:  tag,
		Attributes:  attrs,
	}
	return message.PostMessage(q.config, msg)
}
