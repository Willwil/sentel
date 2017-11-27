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

import "github.com/cloustone/sentel/core"

type EventHandler func(e *Event, ctx interface{})

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
	if deploy, err := c.String("broker", "deployment"); err == nil && deploy == "cluster" {
		return newClusterEventManager(c)
	} else {
		return newLocalEventManager(c)
	}
}
