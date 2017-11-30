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
	"sync"

	eap "github.com/eapache/queue"
)

// Synchronous FIFO queue
type syncQueue struct {
	lock    sync.Mutex
	popable *sync.Cond
	buffer  *eap.Queue
	closed  bool
}

// Create a new syncQueue
func newSyncQueue() *syncQueue {
	ch := &syncQueue{
		buffer: eap.New(),
	}
	ch.popable = sync.NewCond(&ch.lock)
	return ch
}

// Pop an item from syncQueue, will block if syncQueue is empty
func (p *syncQueue) pop() (v interface{}) {
	c := p.popable
	buffer := p.buffer

	p.lock.Lock()
	for buffer.Length() == 0 && !p.closed {
		c.Wait()
	}

	if buffer.Length() > 0 {
		v = buffer.Peek()
		buffer.Remove()
	}

	p.lock.Unlock()
	return
}

func (p *syncQueue) front() interface{} {
	buffer := p.buffer
	p.lock.Lock()
	v := buffer.Peek()
	p.lock.Unlock()
	return v
}

// Try to pop an item from syncQueue, will return immediately with bool=false if syncQueue is empty
func (p *syncQueue) tryPop() (v interface{}, ok bool) {
	buffer := p.buffer

	p.lock.Lock()

	if buffer.Length() > 0 {
		v = buffer.Peek()
		buffer.Remove()
		ok = true
	} else if p.closed {
		ok = true
	}

	p.lock.Unlock()
	return
}

// Push an item to syncQueue. Always returns immediately without blocking
func (p *syncQueue) push(v interface{}) {
	p.lock.Lock()
	if !p.closed {
		p.buffer.Add(v)
		p.popable.Signal()
	}
	p.lock.Unlock()
}

// Get the length of syncQueue
func (p *syncQueue) len() (l int) {
	p.lock.Lock()
	l = p.buffer.Length()
	p.lock.Unlock()
	return
}

// Close syncQueue
// After close, Pop will return nil without block, and TryPop will return v=nil, ok=True
func (p *syncQueue) close() {
	p.lock.Lock()
	if !p.closed {
		p.closed = true
		p.popable.Signal()
	}
	p.lock.Unlock()
}
