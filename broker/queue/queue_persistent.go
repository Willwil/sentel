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
	"time"

	"github.com/cloustone/sentel/core"
)

type persistentQueue struct {
	config   core.Config // Configuration
	id       string      // Queue identifier
	dataChan chan []byte // Data channe
	observer Observer    // Queue observer when data is available
	ref      int         // Ref cound
	plugin   queuePlugin // backend queue plugin
}

func newPersistentQueue(id string, c core.Config) (Queue, error) {
	plugin, err := newPlugin(id, c)
	if err != nil {
		return nil, err
	}
	q := &persistentQueue{
		config:   c,
		id:       id,
		dataChan: make(chan []byte),
		ref:      1,
		plugin:   plugin,
	}
	return q, nil
}

func (p *persistentQueue) Id() string { return p.id }

func (p *persistentQueue) Read(b []byte) (n int, err error) {
	if p.observer != nil {
		data := <-p.dataChan
		return len(data), nil
	}
	data, err := p.plugin.getData()
	if err != nil {
		return -1, err
	}
	b = append(b, data.Data...)
	return len(b), nil
}

// Write writes data to the connection.
// Write can be made to time out and return a Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline.
func (p *persistentQueue) Write(b []byte) (n int, err error) {
	if p.observer == nil {
		// if observer is not set, the data should be write to backend data
		// queue
		p.plugin.pushData(&queueData{Time: time.Now(), Data: b})
	} else {
		p.dataChan <- b
		p.observer.DataAvailable(p, len(b))
	}
	return len(b), nil
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (p *persistentQueue) Close() error {
	close(p.dataChan)
	return nil
}

func (p *persistentQueue) IsPersistent() bool { return false }
func (p *persistentQueue) Name() string       { return p.id }

// RegisterObesrve register an observer on queue
func (p *persistentQueue) RegisterObserver(o Observer) {
	p.observer = o
}

func (p *persistentQueue) Release() int {
	p.ref -= 1
	return p.ref
}
