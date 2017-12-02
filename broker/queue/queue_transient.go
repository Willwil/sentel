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
	"github.com/golang/glog"
)

type transientQueue struct {
	config   core.Config
	clientId string
	observer Observer
	syncq    *syncQueue
}

func newTransientQueue(clientId string, c core.Config, o Observer) (Queue, error) {
	q := &transientQueue{
		config:   c,
		clientId: clientId,
		syncq:    newSyncQueue(),
		observer: o,
	}
	return q, nil
}

// Id return queue identifier
func (p *transientQueue) ClientId() string     { return p.clientId }
func (p *transientQueue) Length() int          { return p.syncq.len() }
func (p *transientQueue) Front() *base.Message { return p.syncq.front().(*base.Message) }
func (p *transientQueue) Pop() *base.Message   { return p.syncq.pop().(*base.Message) }

func (p *transientQueue) Pushback(msg *base.Message) {
	p.syncq.push(msg)
	if p.observer != nil {
		p.observer.DataAvailable(p, msg)
	} else {
		glog.Fatal("queue service: observer is null in transient queue")
	}
}

func (p *transientQueue) Close() error {
	p.syncq.close()
	return nil
}

func (p *transientQueue) IsPersistent() bool { return false }
func (p *transientQueue) Name() string       { return p.clientId }

// RegisterObesrve register an observer on queue
func (p *transientQueue) RegisterObserver(o Observer) {
	p.observer = o
}
