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
	"github.com/cloustone/sentel/pkg/config"
	"github.com/golang/glog"
)

type transientQueue struct {
	config   config.Config
	clientID string
	observer Observer
	syncq    *syncQueue
	list     List
}

func newTransientQueue(clientID string, c config.Config) (Queue, error) {
	q := &transientQueue{
		config:   c,
		clientID: clientID,
		syncq:    newSyncQueue(),
		observer: nil,
		list:     newTransientList(),
	}
	return q, nil
}

// Id return queue identifier
func (p *transientQueue) ClientId() string     { return p.clientID }
func (p *transientQueue) Length() int          { return p.syncq.len() }
func (p *transientQueue) Front() *base.Message { return p.syncq.front().(*base.Message) }
func (p *transientQueue) Pop() *base.Message   { return p.syncq.pop().(*base.Message) }

func (p *transientQueue) Pushback(msg *base.Message) {
	p.syncq.push(msg)
	if p.observer != nil {
		p.observer.DataAvailable(p, msg)
	} else {
		glog.Error("queue service: observer is null in transient queue")
	}
}

func (p *transientQueue) Close() error {
	p.syncq.close()
	return nil
}

func (p *transientQueue) IsPersistent() bool { return false }
func (p *transientQueue) Name() string       { return p.clientID }

// RegisterObesrve register an observer on queue
func (p *transientQueue) RegisterObserver(o Observer) {
	p.observer = o
}

func (p *transientQueue) GetRetainList() List { return p.list }
