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

package executor

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

// ExecutorService manage all rule engine and execute rule
type ExecutorService struct {
	core.ServiceBase
	ruleChan chan *Rule
	engines  map[string]*ruleEngine
	mutex    sync.Mutex
}

type ExecutorServiceFactory struct{}

//  A global executor service instance is needed because indication will send rule to it
var executorService *ExecutorService

// New create executor service factory
func (p *ExecutorServiceFactory) New(c core.Config, quit chan os.Signal) (core.Service, error) {
	// try connect with mongo db
	hosts := c.MustString("conductor", "mongo")
	timeout := c.MustInt("conductor", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	session.Close()

	executorService = &ExecutorService{
		ServiceBase: core.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		ruleChan: make(chan *Rule),
		engines:  make(map[string]*ruleEngine),
		mutex:    sync.Mutex{},
	}
	return executorService, nil
}

// Name
func (p *ExecutorService) Name() string {
	return "executor"
}

// Start
func (p *ExecutorService) Start() error {
	// start rule channel
	go func(s *ExecutorService) {
		s.WaitGroup.Add(1)
		select {
		case r := <-s.ruleChan:
			s.handleRule(r)
		case <-s.Quit:
			break
		}
	}(p)
	return nil
}

// Stop
func (p *ExecutorService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)

	// stop all ruleEngine
	for _, engine := range p.engines {
		if engine != nil {
			engine.stop()
		}
	}
}

// handleRule process incomming rule
func (p *ExecutorService) handleRule(r *Rule) error {
	// Get engine instance according to product id
	if _, ok := p.engines[r.ProductId]; !ok { // not found
		engine, err := newRuleEngine(p.Config, r.ProductId)
		if err != nil {
			glog.Errorf("Failed to create rule engint for product(%s)", r.ProductId)
			return err
		}
		p.engines[r.ProductId] = engine
	}
	engine := p.engines[r.ProductId]

	switch r.RuleAction {
	case core.RuleActionCreate:
		return engine.createRule(r)
	case core.RuleActionRemove:
		return engine.removeRule(r)
	case core.RuleActionUpdate:
		return engine.updateRule(r)
	case core.RuleActionStart:
		return engine.startRule(r)
	case core.RuleActionStop:
		return engine.stopRule(r)
	}
	return nil
}

// HandleRuleNotification handle rule notifications recevied from kafka,
// it will check rule's validity,for example, wether rule exist in database.
func HandleRuleNotification(r *Rule) error {
	glog.Infof("New rule notification: ruleId=%s, ruleName=%s, action=%s", r.RuleName, r.RuleName, r.Action)

	// Check action's validity
	switch r.RuleAction {
	case core.RuleActionCreate:
	case core.RuleActionRemove:
	case core.RuleActionUpdate:
	case core.RuleActionStart:
	case core.RuleActionStop:
	default:
		return fmt.Errorf("Invalid rule action(%s) for product(%s)", r.Action, r.ProductId)
	}

	if r.RuleName == "" || r.ProductId == "" {
		return fmt.Errorf("Invalid argument")
	}
	// Now just simply send rule to executor
	executorService.ruleChan <- r
	return nil
}
