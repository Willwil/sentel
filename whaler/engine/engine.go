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
	"github.com/golang/glog"
)

type RuleEngine interface {
	Start() error
	Stop()
	Recovery() error
	//SetCleanDuration(time.Duration)
	HandleRule(*RuleContext) error
}

func NewEngine(c config.Config) (RuleEngine, error) {
	return newRuleEngine(c)
}

// ruleEngine manage all rule engine and execute rule
type ruleEngine struct {
	config        config.Config            // global configuration
	waitgroup     sync.WaitGroup           // waitigroup for goroutine exit
	quitChan      chan interface{}         // exit channel
	ruleChan      chan *RuleContext        // rule channel to receive rule notification
	executors     map[string]*ruleExecutor // all product's executors
	mutex         sync.Mutex               // mutex
	consumer      message.Consumer         // message consumer for rule creation...
	cleanTimer    *time.Timer              // cleantimer that scan all executors in duration time and clean empty executors
	registry      *registry.Registry
	cleanDuration time.Duration
}

const CLEAN_DURATION = time.Minute * 30

// New create rule engine service factory
func newRuleEngine(c config.Config) (RuleEngine, error) {
	r, err := registry.New(c)
	if err != nil {
		return nil, err
	}

	consumer, _ := message.NewConsumer(c, "whaler-engine")
	return &ruleEngine{
		config:        c,
		waitgroup:     sync.WaitGroup{},
		quitChan:      make(chan interface{}),
		ruleChan:      make(chan *RuleContext),
		executors:     make(map[string]*ruleExecutor),
		mutex:         sync.Mutex{},
		consumer:      consumer,
		registry:      r,
		cleanDuration: CLEAN_DURATION,
	}, nil
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
				err := s.dispatchRule(ctx)
				if ctx.Response != nil {
					ctx.Response <- err
				}
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
	p.cleanTimer.Stop()
	p.quitChan <- true
	p.waitgroup.Wait()
	close(p.quitChan)
	close(p.ruleChan)
	p.consumer.Close()

	// stop all ruleEngine
	for _, executor := range p.executors {
		if executor != nil {
			executor.stop()
		}
	}
}

// Recovery restart rules after whaler is restarted
func (p *ruleEngine) Recovery() error {
	rules := p.registry.GetRulesWithStatus(registry.RuleStatusStarted)
	for _, r := range rules {
		if r.Status == registry.RuleStatusStarted {
			ctx := NewRuleContext(r.ProductId, r.RuleName, message.ActionCreate)
			if err := p.HandleRule(ctx); err != nil {
				glog.Errorf("product '%s', rule '%s'recovery create failed", r.ProductId, r.RuleName)
			}
			ctx.Action = message.ActionStart
			if err := p.HandleRule(ctx); err != nil {
				glog.Errorf("product '%s', rule '%s'recovery start failed", r.ProductId, r.RuleName)
			}

		}
	}
	return nil
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

// dispatchRule process incomming rule according to rule's action
func (p *ruleEngine) dispatchRule(ctx *RuleContext) error {
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

func (p *ruleEngine) messageHandlerFunc(msg message.Message, param interface{}) {
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
		ctx := NewRuleContext(r.ProductId, r.RuleName, action)
		p.ruleChan <- ctx
		return
	}
	glog.Errorf("invalid message topic '%s' with rule", msg.Topic())
}

func (p *ruleEngine) HandleRule(ctx *RuleContext) error {
	ctx.Response = make(chan error)
	p.ruleChan <- ctx
	err := <-ctx.Response
	return err
}
