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

package event

import (
	"os"
	"sync"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/core"
)

type eventService struct {
	base.ServiceBase
	eventmgr eventManager
}

const (
	ServiceName = "event"
)

func New(c core.Config, quit chan os.Signal) (base.Service, error) {
	eventmgr, err := newEventManager(c)
	if err != nil {
		return nil, err
	}
	return &eventService{
		ServiceBase: base.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		eventmgr: eventmgr,
	}, nil

}

// Name
func (p *eventService) Name() string {
	return ServiceName
}

func (p *eventService) Initialize() error {
	return p.eventmgr.initialize(p.Config)
}

// Start
func (p *eventService) Start() error {
	return p.eventmgr.run()
}

// Stop
func (p *eventService) Stop() {
	p.eventmgr.stop()
}

func (p *eventService) subscribe(event uint32, handler EventHandler, ctx interface{}) {
	p.eventmgr.subscribe(event, handler, ctx)
}
func (p *eventService) notify(event *Event) {
	p.eventmgr.notify(event)
}
