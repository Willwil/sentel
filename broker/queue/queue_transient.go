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
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

type transientQueue struct {
	config   core.Config
	id       string
	dataChan chan []byte
	observer Observer
}

func newTransientQueue(id string, c core.Config) (Queue, error) {
	q := &transientQueue{
		config:   c,
		id:       id,
		dataChan: make(chan []byte),
	}
	return q, nil
}

// Id return queue identifier
func (p *transientQueue) Id() string { return p.id }

// Read read bytes from queue
func (p *transientQueue) Read(b []byte) (n int, err error) {
	data := <-p.dataChan
	return len(data), nil
}

// Write writes data to the connection.
// Write can be made to time out and return a Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline.
func (p *transientQueue) Write(b []byte) (n int, err error) {
	p.dataChan <- b
	if p.observer != nil {
		p.observer.DataAvailable(p, len(b))
	} else {
		glog.Fatal("queue service: observer is null in transient queue")
	}
	return len(b), nil
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (p *transientQueue) Close() error {
	close(p.dataChan)
	return nil
}

func (p *transientQueue) IsPersistent() bool { return false }
func (p *transientQueue) Name() string       { return p.id }

// RegisterObesrve register an observer on queue
func (p *transientQueue) RegisterObserver(o Observer) {
	p.observer = o
}
