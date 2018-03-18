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

package mns

import (
	"sync"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
)

type manageService struct {
	config    config.Config
	waitgroup sync.WaitGroup
}
type ServiceFactory struct{}

func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	return &manageService{
		config:    c,
		waitgroup: sync.WaitGroup{},
	}, nil
}

func (p *manageService) Name() string      { return "management" }
func (p *manageService) Initialize() error { return nil }

func (p *manageService) Start() error {
	p.waitgroup.Add(1)
	go func(s *manageService) {
		defer s.waitgroup.Done()
	}(p)
	return nil
}

func (p *manageService) Stop() {
	p.waitgroup.Wait()
}
