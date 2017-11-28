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

package event

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cloustone/sentel/core"
)

type localEventManager struct {
	eventChan   chan *Event                    // Event notify channel
	subscribers map[uint32][]subscriberContext // All subscriber's contex
	mutex       sync.Mutex                     // Mutex to lock subscribers
	quit        chan os.Signal
	waitgroup   sync.WaitGroup
	product     string
}

const (
	localEventManagerVersion = "0.1"
)

// newlocalEventManager create global broker
func newLocalEventManager(product string, c core.Config) (*localEventManager, error) {
	return &localEventManager{
		eventChan:   make(chan *Event),
		subscribers: make(map[uint32][]subscriberContext),
		quit:        make(chan os.Signal),
		waitgroup:   sync.WaitGroup{},
		product:     product,
	}, nil
}

// initialize
func (p *localEventManager) initialize(c core.Config) error {
	return nil
}

// Start
func (p *localEventManager) run() error {
	go func(p *localEventManager) {
		for {
			select {
			case e := <-p.eventChan:
				// When event is received from local service, we should
				// transeverse other service to process it at first
				if _, found := p.subscribers[e.Type]; found {
					subscribers := p.subscribers[e.Type]
					for _, subscriber := range subscribers {
						go subscriber.handler(e, subscriber.ctx)
					}
				}
			case <-p.quit:
				return
			}
		}
	}(p)
	return nil
}

// Stop
func (p *localEventManager) stop() {
	signal.Notify(p.quit, syscall.SIGINT, syscall.SIGQUIT)
	close(p.eventChan)
}

// publish publish event to event service
func (p *localEventManager) notify(e *Event) {
	p.eventChan <- e
}

// subscribe subcribe event from event service
func (p *localEventManager) subscribe(t uint32, handler EventHandler, ctx interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.subscribers[t]; !ok {
		// Create handler list if not existed
		p.subscribers[t] = []subscriberContext{}
	}
	p.subscribers[t] = append(p.subscribers[t], subscriberContext{handler: handler, ctx: ctx})
}
