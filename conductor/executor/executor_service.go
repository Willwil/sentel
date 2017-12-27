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

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

// executorService manage all rule engine and execute rule
type executorService struct {
	service.ServiceBase
	ruleChan chan ruleContext
	engines  map[string]*ruleEngine
	mutex    sync.Mutex
}

type ruleContext struct {
	rule   *registry.Rule
	action string
}

type ServiceFactory struct{}

// New create executor service factory
func (p ServiceFactory) New(c config.Config, quit chan os.Signal) (service.Service, error) {
	// try connect with mongo db
	hosts := c.MustString("conductor", "mongo")
	timeout := c.MustInt("conductor", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	session.Close()

	return &executorService{
		ServiceBase: service.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		ruleChan: make(chan ruleContext),
		engines:  make(map[string]*ruleEngine),
		mutex:    sync.Mutex{},
	}, nil
}

// Name
func (p *executorService) Name() string { return "executor" }

// Start
func (p *executorService) Start() error {
	// start rule channel
	go func(s *executorService) {
		s.WaitGroup.Add(1)
		select {
		case ctx := <-s.ruleChan:
			s.handleRule(ctx)
		case <-s.Quit:
			break
		}
	}(p)
	return nil
}

// Stop
func (p *executorService) Stop() {
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
func (p *executorService) handleRule(ctx ruleContext) error {
	r := ctx.rule
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

	switch ctx.action {
	case message.RuleActionCreate:
		return engine.createRule(r)
	case message.RuleActionRemove:
		return engine.removeRule(r)
	case message.RuleActionUpdate:
		return engine.updateRule(r)
	case message.RuleActionStart:
		return engine.startRule(r)
	case message.RuleActionStop:
		return engine.stopRule(r)
	}
	return nil
}

// HandleRuleNotification handle rule notifications recevied from kafka,
// it will check rule's validity,for example, wether rule exist in database.
func HandleRuleNotification(t *message.RuleTopic) error {
	glog.Infof("New rule notification for  rule '%s'", t.RuleName)

	// Check action's validity
	switch t.RuleAction {
	case message.RuleActionCreate:
	case message.RuleActionRemove:
	case message.RuleActionUpdate:
	case message.RuleActionStart:
	case message.RuleActionStop:
	default:
		return fmt.Errorf("Invalid rule action(%s) for product(%s)", t.RuleAction, t.ProductId)
	}
	// Now just simply send rule to executor
	executor := service.GetServiceManager().GetService("executor").(*executorService)
	ctx := ruleContext{
		action: t.RuleAction,
		rule: &registry.Rule{
			ProductId: t.ProductId,
			RuleName:  t.RuleName,
		},
	}
	executor.ruleChan <- ctx
	return nil
}
