//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package iothub

import (
	"errors"
	"sync"

	mgo "gopkg.in/mgo.v2"

	"github.com/cloustone/sentel/core"
)

type NotifyService struct {
	config core.Config
	chn    chan core.ServiceCommand
	wg     sync.WaitGroup
}

// NotifyServiceFactory
type NotifyServiceFactory struct{}

// New create apiService service factory
func (m *NotifyServiceFactory) New(protocol string, c core.Config, ch chan core.ServiceCommand) (core.Service, error) {
	// check mongo db configuration
	hosts, err := c.String("iothub", "mongo")
	if err != nil || hosts == "" {
		return nil, errors.New("Invalid mongo configuration")
	}

	// try connect with mongo db
	session, err := mgo.Dial(hosts)
	if err != nil {
		return nil, err
	}
	session.Close()

	return &NotifyService{
		config: c,
		wg:     sync.WaitGroup{},
		chn:    ch,
	}, nil
}

// Name
func (s *NotifyService) Name() string {
	return "notify-service"
}

// Start
func (s *NotifyService) Start() error {
	go func(s *NotifyService) {
		s.wg.Add(1)
	}(s)
	return nil
}

// Stop
func (s *NotifyService) Stop() {
	s.wg.Wait()
}