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

package broker

import (
	"errors"

	"github.com/cloustone/sentel/core"
)

type EventHandler func(e *Event, ctx interface{})

// eventManager manage all events shared by brokers
type eventManager interface {
	initialize(c core.Config) error
	run() error
	stop()
	subscribe(event uint32, handler EventHandler, ctx interface{})
	notify(event *Event)
}

// subscriberContext hold subscriber handler and context
type subscriberContext struct {
	handler EventHandler
	ctx     interface{}
}

// newEventManger create event manager according to deploy mode
func newEventManager(c core.Config) (eventManager, error) {
	// get deployment mode and product name
	deploy, dErr := c.String("broker", "deploy")
	product, pErr := c.String("broker", "product")

	if dErr != nil || pErr != nil {
		return nil, errors.New("Invalid event manager options")
	}

	if deploy == "cluster" {
		return newClusterEventManager(product, c)
	} else {
		return newLocalEventManager(product, c)
	}
}
