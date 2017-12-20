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
	"github.com/cloustone/sentel/common"
)

type persistentQueue struct {
	config   com.Config  // Configuration
	clientId string      // Queue identifier
	observer Observer    // Queue observer when data is available
	plugin   queuePlugin // backend queue plugin
}

func newPersistentQueue(clientId string, c com.Config, o Observer) (Queue, error) {
	plugin, err := newPlugin(clientId, c)
	if err != nil {
		return nil, err
	}
	q := &persistentQueue{
		config:   c,
		clientId: clientId,
		plugin:   plugin,
		observer: o,
	}
	return q, nil
}

func (p *persistentQueue) ClientId() string     { return p.clientId }
func (p *persistentQueue) Length() int          { return p.plugin.length() }
func (p *persistentQueue) Front() *base.Message { return p.plugin.front() }

func (p *persistentQueue) Pushback(msg *base.Message) {
	p.plugin.pushback(msg)
	if p.observer != nil {
		p.observer.DataAvailable(p, msg)
	}
}

func (p *persistentQueue) Pop() *base.Message          { return p.plugin.pop() }
func (p *persistentQueue) Close() error                { return p.plugin.close() }
func (p *persistentQueue) IsPersistent() bool          { return true }
func (p *persistentQueue) Name() string                { return p.clientId }
func (p *persistentQueue) RegisterObserver(o Observer) { p.observer = o }
