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
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
)

type transientQueue struct {
	id       string
	conn     net.Conn
	dataChan chan []byte
	quit     chan os.Signal
}

func newTransientQueue(id string, conn net.Conn) (Queue, error) {
	q := &transientQueue{
		id:       id,
		conn:     conn,
		dataChan: make(chan []byte),
		quit:     make(chan os.Signal),
	}
	go func(q *transientQueue) {
		for {
			select {
			case data := <-q.dataChan:
				wsize := 0
				for wsize < len(data) {
					len, err := q.conn.Write(data[wsize:])
					if err != nil {
						glog.Fatal("broker: queue '%s' write failed", q.id)
						continue
					}
					wsize += len
				}
			case <-q.quit:
				return
			}
		}
	}(q)
	return q, nil
}

func (p *transientQueue) Read(b []byte) (n int, err error) {
	return 0, nil
}

// Write writes data to the connection.
// Write can be made to time out and return a Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline.
func (p *transientQueue) Write(b []byte) (n int, err error) {
	p.dataChan <- b
	return len(b), nil
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (p *transientQueue) Close() error {
	signal.Notify(p.quit, syscall.SIGINT, syscall.SIGQUIT)
	close(p.quit)
	close(p.dataChan)
	return nil
}

func (p *transientQueue) IsPersistent() bool { return false }
func (p *transientQueue) Name() string       { return p.id }
