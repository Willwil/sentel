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

type Observer interface {
	// DataAvailable notify that data is available
	DataAvailable(q Queue, n int)
}

type Queue interface {
	// Id return identifier of queue
	Id() string

	// RegisterObesrve register an observer on queue
	RegisterObserver(o Observer)

	// Read reads data from the queue.
	Read(b []byte) (n int, err error)

	// Write writes data to the queue.
	Write(b []byte) (n int, err error)

	// Close closes the connection.
	// Any blocked Read or Write operations will be unblocked and return errors.
	Close() error

	// IsPersistent return true if queue is persistent
	IsPersistent() bool
}
