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

const minQueueLen = 16

type transientList struct {
	buf   []*base.Message
	count int
}

func newTransientList() List {
	return &transientList{
		buf: make([]*base.Message, minQueueLen),
	}
}

func (p *transientList) resize() {
	newBuf := make([]*base.Message, p.count<<1)
	copy(newBuf, p.buf[:p.count])
	p.buf = newBuf
}

// Length return total mesasge count in List
func (p *transientList) Length() int {
	return p.count
}

// Pushback push message at tail of List
func (p *transientList) Pushback(msg *base.Message) {
	if p.count == len(p.buf) {
		p.resize()
	}
	p.buf[p.count] = msg
	p.count++
}

// Get returns the element at index i in the List
func (p *transientList) Get(index int) *base.Message {
	if index >= p.count {
		return nil
	}
	return p.buf[index]
}

// Remove
func (p *transientList) Remove(pktID uint16) {
	for i := 0; i < p.count; i++ {
		if p.buf[i].PacketId == pktID {
			p.buf = append(p.buf[:i], p.buf[i+1:]...)
			p.count--
			break
		}
	}
	if len(p.buf) > minQueueLen && (p.count<<2) == len(p.buf) {
		p.resize()
	}
}
