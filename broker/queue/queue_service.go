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
	"fmt"
	"os"
	"sync"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/pkg/config"
)

// Metaservice manage broker metadata
// Broker's metadata include the following data
// - Global broker cluster data
// - Shadow device
type queueService struct {
	base.ServiceBase
	queues map[string]Queue // All queues
	mutex  sync.Mutex       // Mutex for queues
}

const (
	ServiceName = "queue"
)

type ServiceFactory struct{}

// New create metadata service factory
func (p ServiceFactory) New(c config.Config, quit chan os.Signal) (base.Service, error) {
	return &queueService{
		ServiceBase: base.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		queues: make(map[string]Queue),
		mutex:  sync.Mutex{},
	}, nil

}

// Name
func (p *queueService) Name() string {
	return ServiceName
}

func (p *queueService) Initialize() error {
	return p.bootstrap()
}

func (p *queueService) bootstrap() error {
	return nil
}

// Start
func (p *queueService) Start() error {
	go func(p *queueService) {
		for {
			select {
			case <-p.Quit:
				return
			}
		}
	}(p)
	return nil
}

// Stop
func (p *queueService) Stop() {
	p.WaitGroup.Wait()
}

// newQueue allocate queue from queue service
func (p *queueService) newQueue(id string, persistent bool, o Observer) (Queue, error) {
	p.mutex.Lock()
	p.mutex.Unlock()

	// check wethere the id's queue already existed
	if q, found := p.queues[id]; found {
		if persistent && q.IsPersistent() {
			return q, nil
		}

		if q.IsPersistent() && !persistent {
			p.destroyQueue(id)
		} else {
			return nil, fmt.Errorf("broker: queue for '%s' already exist", id)
		}
	}

	var q Queue
	var err error
	if persistent {
		q, err = newPersistentQueue(id, p.Config, o)
	} else {
		q, err = newTransientQueue(id, p.Config, o)
	}
	if err != nil {
		return nil, err
	}
	p.queues[id] = q
	return q, nil
}

// freeQueue release queue from queue service
func (p *queueService) destroyQueue(id string) {
	p.mutex.Lock()
	p.mutex.Unlock()
	delete(p.queues, id)
}

// getQueue return queue by id
func (p *queueService) getQueue(id string) Queue {
	if _, found := p.queues[id]; found {
		return p.queues[id]
	}
	return nil
}

// releaseQueue decrease queue's reference count, and destory the queue if reference is zero
func (p *queueService) releaseQueue(id string) {
	if _, found := p.queues[id]; found {
		p.destroyQueue(id)
	}
}
