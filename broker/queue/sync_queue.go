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
func (q *syncQueue) pop() (v interface{}) {
	c := q.popable
	buffer := q.buffer

	q.lock.Lock()
	for buffer.Length() == 0 && !q.closed {
		c.Wait()
	}

	if buffer.Length() > 0 {
		v = buffer.Peek()
		buffer.Remove()
	}

	q.lock.Unlock()
	return
}

func (q *syncQueue) front() interface{} {
	buffer := q.buffer
	q.lock.Lock()
	v := buffer.Peek()
	q.lock.Unlock()
	return v
}

// Try to pop an item from syncQueue, will return immediately with bool=false if syncQueue is empty
func (q *syncQueue) tryPop() (v interface{}, ok bool) {
	buffer := q.buffer

	q.lock.Lock()

	if buffer.Length() > 0 {
		v = buffer.Peek()
		buffer.Remove()
		ok = true
	} else if q.closed {
		ok = true
	}

	q.lock.Unlock()
	return
}

// Push an item to syncQueue. Always returns immediately without blocking
func (q *syncQueue) push(v interface{}) {
	q.lock.Lock()
	if !q.closed {
		q.buffer.Add(v)
		q.popable.Signal()
	}
	q.lock.Unlock()
}

// Get the length of syncQueue
func (q *syncQueue) len() (l int) {
	q.lock.Lock()
	l = q.buffer.Length()
	q.lock.Unlock()
	return
}

// Close syncQueue
// After close, Pop will return nil without block, and TryPop will return v=nil, ok=True
func (q *syncQueue) close() {
	q.lock.Lock()
	if !q.closed {
		q.closed = true
		q.popable.Signal()
	}
	q.lock.Unlock()
}
