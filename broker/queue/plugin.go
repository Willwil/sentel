//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package queue

import (
	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/core"
)

type queuePlugin interface {
	// Initialize initialize the backend queue plugin
	initialize() error
	// GetMessageCount return total mesasge count in queue
	length() int

	// Front return message at queue's head
	front() *base.Message

	// Pushback push message at tail of queue
	pushback(msg *base.Message)

	// Pop popup head message from queue
	pop() *base.Message

	// Close closes the connection.
	close() error
}

// newPlugin return backend queue plugin accoriding to configuration and id
func newPlugin(clientId string, c core.Config) (queuePlugin, error) {
	q, err := newMongoPlugin(clientId, c)
	if q != nil {
		err = q.initialize()
		return q, err
	}
	return nil, err
}
