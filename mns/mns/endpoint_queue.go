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

// queueEndpoint has scheme "mns:AccountId:queues/queueName"
type queueEndpoint struct {
	accountId string
	queueName string
	attribute EndpointAttribute
}

func newQueueEndpoint(c config.Config, uri string) (endpoint Endpoint, err error) {
	accountId, queueName, err := parseQueueScheme(uri)
	endpoint = &queueEndpoint{
		accountId: accountId,
		queueName: queueName,
		attribute: EndpointAttribute{
			Name:           uri,
			Type:           "mns",
			URI:            uri,
			CreatedAt:      time.Now(),
			LastModifiedAt: time.Now(),
		},
	}
	return
}

func parseQueueScheme(uri string) (string, string, error) {
	accountId := ""
	queueName := ""
	err := ErrInvalidArgument
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
func (q queueEndpoint) PushMessage(msg Message) error   { return nil }
