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

	"github.com/cloustone/sentel/common"
	"github.com/cloustone/sentel/common/db"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

// executorService manage all rule engine and execute rule
type executorService struct {
	com.ServiceBase
	ruleChan chan ruleContext
	engines  map[string]*ruleEngine
	mutex    sync.Mutex
}

type ruleContext struct {
	rule   *db.Rule
	action string
}

type ServiceFactory struct{}

// New create executor service factory
func (p ServiceFactory) New(c com.Config, quit chan os.Signal) (com.Service, error) {
	// try connect with mongo db
	hosts := c.MustString("conductor", "mongo")
	timeout := c.MustInt("conductor", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	session.Close()

	return &executorService{
		ServiceBase: com.ServiceBase{
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
		engine, err := newRuleEngine(p.Config, r.TenantId, r.ProductId)
		if err != nil {
			glog.Errorf("Failed to create rule engint for product(%s)", r.ProductId)
			return err
		}
		p.engines[r.ProductId] = engine
	}
	engine := p.engines[r.ProductId]

	switch ctx.action {
	case com.RuleActionCreate:
		return engine.createRule(r)
	case com.RuleActionRemove:
		return engine.removeRule(r)
	case com.RuleActionUpdate:
		return engine.updateRule(r)
	case com.RuleActionStart:
		return engine.startRule(r)
	case com.RuleActionStop:
		return engine.stopRule(r)
	}
	return nil
}

// HandleRuleNotification handle rule notifications recevied from kafka,
// it will check rule's validity,for example, wether rule exist in database.
func HandleRuleNotification(t *com.RuleTopic) error {
	glog.Infof("New rule notification for  rule '%s'", t.RuleName)

	// Check action's validity
	switch t.RuleAction {
	case com.RuleActionCreate:
	case com.RuleActionRemove:
	case com.RuleActionUpdate:
	case com.RuleActionStart:
	case com.RuleActionStop:
	default:
		return fmt.Errorf("Invalid rule action(%s) for product(%s)", t.RuleAction, t.ProductId)
	}
	// Now just simply send rule to executor
	executor := com.GetServiceManager().GetService("executor").(*executorService)
	ctx := ruleContext{
		action: t.RuleAction,
		rule: &db.Rule{
			ProductId: t.ProductId,
			RuleName:  t.RuleName,
		},
	}
	executor.ruleChan <- ctx
	return nil
}
