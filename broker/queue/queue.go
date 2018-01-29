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

import "github.com/cloustone/sentel/broker/base"

type Observer interface {
	// DataAvailable notify that data is available
	DataAvailable(q Queue, msg *base.Message)
}

type List interface {
	// Length return total mesasge count in List
	Length() int

	// Pushback push message at tail of List
	Pushback(msg *base.Message)

	// Get returns the element at index i in the List
	Get(index int) *base.Message

	// Remove
	Remove(pktID uint16)
}

type Queue interface {
	// Id return client id of queue
	ClientId() string

	// RegisterObesrve register an observer on queue
	RegisterObserver(o Observer)

	// GetMessageCount return total mesasge count in queue
	Length() int

	// Front return message at queue's head
	Front() *base.Message

	// Pushback push message at tail of queue
	Pushback(msg *base.Message)

	// Pop popup head message from queue
	Pop() *base.Message

	// GetRetainList get the message list waiting for PUBACK
	GetRetainList() List

	// Close closes the connection.
	Close() error

	// IsPersistent return true if queue is persistent
	IsPersistent() bool
}
