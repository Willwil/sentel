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

package engine

import (
	"sync"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/cloustone/sentel/pkg/service"
)

// ruleEngine manage all rule engine and execute rule
type engineService struct {
	config       config.Config  // global configuration
	waitgroup    sync.WaitGroup // waitigroup for goroutine exit
	recoveryChan chan bool      // recover channle that load started rule at the startup timeing
	ruleEngine   RuleEngine
}

const SERVICE_NAME = "engine"

type ServiceFactory struct{}

// New create rule engine service factory
func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	ruleEngine, err := NewEngine(c)
	if err != nil {
		return nil, err
	}

	return &engineService{
		config:       c,
		waitgroup:    sync.WaitGroup{},
		recoveryChan: make(chan bool, 1),
		ruleEngine:   ruleEngine,
	}, nil
}

// Name
func (p *engineService) Name() string { return SERVICE_NAME }

// Initialize check all rule's state and recover if rule is started
func (p *engineService) Initialize() error {
	if r, err := registry.New(p.config); err == nil {
		defer r.Close()
		rules := r.GetRulesWithStatus(registry.RuleStatusStarted)
		p.recoveryChan <- (len(rules) > 0)
	}
	return nil
}

// Start
func (p *engineService) Start() error {
	if err := p.ruleEngine.Start(); err != nil {
		return err
	}
	p.waitgroup.Add(1)
	go func(s *engineService) {
		defer s.waitgroup.Done()
		for {
			select {
			case recovery := <-s.recoveryChan:
				if recovery {
					s.ruleEngine.Recovery()
				}
				return
			}
		}
	}(p)
	return nil
}

// Stop
func (p *engineService) Stop() {
	p.recoveryChan <- false
	p.waitgroup.Wait()
	close(p.recoveryChan)
	p.ruleEngine.Stop()
}

func HandleRuleNotification(ctx *RuleContext) error {
	mgr := service.GetServiceManager()
	s := mgr.GetService(SERVICE_NAME).(*engineService)
	return s.ruleEngine.HandleRule(ctx)
}
