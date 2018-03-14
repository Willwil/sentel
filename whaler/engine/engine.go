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
	"fmt"
	"sync"
	"time"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"
)

// ruleEngine manage all rule engine and execute rule
type ruleEngine struct {
	config       config.Config            // global configuration
	waitgroup    sync.WaitGroup           // waitigroup for goroutine exit
	quitChan     chan interface{}         // exit channel
	ruleChan     chan RuleContext         // rule channel to receive rule notification
	executors    map[string]*ruleExecutor // all product's executors
	mutex        sync.Mutex               // mutex
	consumer     message.Consumer         // message consumer for rule creation...
	recoveryChan chan interface{}         // recover channle that load started rule at the startup timeing
	cleanTimer   *time.Timer              // cleantimer that scan all executors in duration time and clean empty executors
}

const SERVICE_NAME = "engine"
const CLEAN_DURATION = time.Minute * 30

type ServiceFactory struct{}

// New create rule engine service factory
func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	consumer, _ := message.NewConsumer(c, "whaler-engine")
	return &ruleEngine{
		config:       c,
		waitgroup:    sync.WaitGroup{},
		quitChan:     make(chan interface{}),
		ruleChan:     make(chan RuleContext),
		executors:    make(map[string]*ruleExecutor),
		mutex:        sync.Mutex{},
		consumer:     consumer,
		recoveryChan: make(chan interface{}),
	}, nil
}

// Name
func (p *ruleEngine) Name() string { return SERVICE_NAME }

// Initialize check all rule's state and recover if rule is started
func (p *ruleEngine) Initialize() error {
	// try connect with registry
	r, err := registry.New(p.config)
	if err != nil {
		return err
	}
	defer r.Close()
	// TODO: opening registry twice and read all rules into service is not good here
	rules := r.GetRulesWithStatus(registry.RuleStatusStarted)
	if len(rules) > 0 {
		p.recoveryChan <- true
	}
	return nil
}

// Start
func (p *ruleEngine) Start() error {
	if err := p.consumer.Subscribe(message.TopicNameRule, p.messageHandlerFunc, nil); err != nil {
		return fmt.Errorf("subscribe message failed : %s", err.Error())
	}
	p.consumer.Start()
	p.waitgroup.Add(1)
	p.cleanTimer = time.NewTimer(CLEAN_DURATION)

	go func(s *ruleEngine) {
		defer s.waitgroup.Done()
		for {
			select {
			case ctx := <-s.ruleChan:
				err := s.handleRule(ctx)
				if ctx.Resp != nil {
					ctx.Resp <- err
				}
			case <-s.recoveryChan:
				s.recovery()
			case <-s.quitChan:
				return
			case <-s.cleanTimer.C:
				p.cleanExecutors()
				s.cleanTimer.Reset(CLEAN_DURATION)
			}
		}
	}(p)
	return nil
}

// Stop
func (p *ruleEngine) Stop() {
	if p.cleanTimer != nil {
		p.cleanTimer.Stop()
	}
	p.quitChan <- true
	p.waitgroup.Wait()
	close(p.quitChan)
	close(p.ruleChan)
	close(p.recoveryChan)
	p.consumer.Close()

	// stop all ruleEngine
	for _, executor := range p.executors {
		if executor != nil {
			executor.stop()
		}
	}
}

// recovery restart rules after whaler is restarted
func (p *ruleEngine) recovery() {
	r, _ := registry.New(p.config)
	defer r.Close()
	rules := r.GetRulesWithStatus(registry.RuleStatusStarted)
	for _, r := range rules {
		if r.Status == registry.RuleStatusStarted {
			ctx := RuleContext{
				Action:    message.ActionCreate,
				ProductId: r.ProductId,
				RuleName:  r.RuleName,
			}
			if err := p.handleRule(ctx); err != nil {
				glog.Errorf("product '%s', rule '%s'recovery create failed", r.ProductId, r.RuleName)
			}
			ctx.Action = message.ActionStart
			if err := p.handleRule(ctx); err != nil {
				glog.Errorf("product '%s', rule '%s'recovery start failed", r.ProductId, r.RuleName)
			}

		}
	}
}

// cleanExecutors scan all executors and stop the executor if
// there are no rules to be started for this executor
func (p *ruleEngine) cleanExecutors() {
	for productId, executor := range p.executors {
		if len(executor.rules) == 0 {
			executor.stop()
			delete(p.executors, productId)
			glog.Infof("executor '%s' deleted", productId)
		}
	}
}

// handleRule process incomming rule
func (p *ruleEngine) handleRule(ctx RuleContext) error {
	glog.Infof("handing rule... %s", ctx.String())

	productId := ctx.ProductId
	switch ctx.Action {
	case RuleActionCreate:
		if _, found := p.executors[productId]; !found {
			if executor, err := newRuleExecutor(p.config, productId); err != nil {
				return err
			} else {
				if err := executor.start(); err != nil {
					return err
				}
				p.executors[productId] = executor
			}
		}
		return p.executors[productId].createRule(ctx)

	case RuleActionRemove:
		if executor, found := p.executors[productId]; found {
			return executor.removeRule(ctx)
		}
	case RuleActionUpdate:
		if executor, found := p.executors[productId]; found {
			return executor.updateRule(ctx)
		}

	case RuleActionStart:
		if executor, found := p.executors[productId]; found {
			return executor.startRule(ctx)
		}

	case RuleActionStop:
		if executor, found := p.executors[productId]; found {
			return executor.stopRule(ctx)
		}
	}
	return fmt.Errorf("invalid operation on product '%s' rule '%s'", productId, ctx.RuleName)
}

func (p *ruleEngine) messageHandlerFunc(msg message.Message, ctx interface{}) {
	if r, ok := msg.(*message.Rule); ok || r != nil {
		action := ""
		switch r.Action {
		case message.ActionCreate:
			action = RuleActionCreate
		case message.ActionRemove:
			action = RuleActionRemove
		case message.ActionUpdate:
			action = RuleActionUpdate
		case message.ActionStart:
			action = RuleActionStart
		case message.ActionStop:
			action = RuleActionStop
		default:
			glog.Errorf("invalid rule action '%s' for product '%s'", r.Action, r.ProductId)
		}

		glog.Infof("received topic '%s' with action '%s'...", msg.Topic(), action)
		rc := RuleContext{
			Action:    action,
			ProductId: r.ProductId,
			RuleName:  r.RuleName,
		}
		p.ruleChan <- rc
		return
	}
	glog.Errorf("invalid message topic '%s' with rule", msg.Topic())
}

func HandleRuleNotification(ctx RuleContext) error {
	mgr := service.GetServiceManager()
	s := mgr.GetService(SERVICE_NAME).(*ruleEngine)
	s.ruleChan <- ctx
	return <-ctx.Resp
}

func HandleTopicNotification(productId string, msg message.Message) error {
	mgr := service.GetServiceManager()
	s := mgr.GetService(SERVICE_NAME).(*ruleEngine)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if executor, found := s.executors[productId]; found {
		executor.messageHandlerFunc(msg, nil)
		return nil
	}
	return fmt.Errorf("invalid product '%s'", productId)
}
