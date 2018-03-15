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

package engine

import (
	"fmt"
	"sync"

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/golang/glog"
)

type ruleExecutor struct {
	config    config.Config                 // Configuration
	tenantId  string                        // Tenant
	productId string                        // One product have one rule engine
	rules     map[string]*rule              // All product's rule
	mutex     sync.Mutex                    // Mutex to protext rules list
	started   bool                          // Indicate wether engined is started
	consumer  message.Consumer              // Consumer for tenant broker event topic
	dataChan  chan *event.TopicPublishEvent // Data channel to exchange broker publish event
	quitChan  chan interface{}              // Quit channel to exit executor safely
	waitgroup sync.WaitGroup                // Waitgroup for channel exition
}

// newRuleExecutor create a rule executor according to product id and configuration
func newRuleExecutor(c config.Config, productId string) (*ruleExecutor, error) {
	// get product's tenant id
	r, err := registry.New(c)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	product, err := r.GetProduct(productId)
	if err != nil {
		return nil, fmt.Errorf("invalid product id '%s'", productId)
	}
	clientId := fmt.Sprintf("whaler-rule-executor-%s", productId)
	consumer, _ := message.NewConsumer(c, clientId)
	return &ruleExecutor{
		tenantId:  product.TenantId,
		productId: productId,
		config:    c,
		rules:     make(map[string]*rule),
		mutex:     sync.Mutex{},
		started:   false,
		consumer:  consumer,
		dataChan:  make(chan *event.TopicPublishEvent, 1),
		quitChan:  make(chan interface{}, 1),
		waitgroup: sync.WaitGroup{},
	}, nil
}

// start will start the rule engine, receiving topic and rule
func (p *ruleExecutor) start() error {
	if !p.started {
		p.waitgroup.Add(1)
		go func(p *ruleExecutor) {
			defer p.waitgroup.Done()
			for {
				select {
				case <-p.quitChan:
					return
				case e := <-p.dataChan:
					p.execute(e)
				}
			}
		}(p)

		topic := fmt.Sprintf(event.FmtOfBrokerEventBus, p.tenantId)
		if err := p.consumer.Subscribe(topic, p.messageHandlerFunc, nil); err != nil {
			return err
		}
		p.consumer.Start()
		p.started = true
	}
	return nil
}

// stop will stop the engine
func (p *ruleExecutor) stop() {
	p.consumer.Close()
	p.quitChan <- true
	p.waitgroup.Wait()
	p.started = false
}

// messageHandlerFunc handle mqtt event from other service
func (p *ruleExecutor) messageHandlerFunc(msg message.Message, ctx interface{}) {
	glog.Infof("executor received notification '%s'", msg.Topic())
	if t, err := event.Decode(msg, event.JSONCodec); err == nil && t.GetType() == event.TopicPublish {
		pubevent := t.(*event.TopicPublishEvent)
		if pubevent.ProductID == p.productId {
			// we can call p.execute(t) here, but in consider of avoiding to block message receiver
			// we use datachannel
			p.dataChan <- pubevent
		}
	}
}

// getRule return rule according to rule context
func (p *ruleExecutor) getRule(ctx *RuleContext) (*rule, error) {
	if r, err := registry.New(p.config); err == nil {
		if rr, err := r.GetRule(ctx.ProductId, ctx.RuleName); err == nil {
			return &rule{Rule: *rr}, nil
		}
	}
	return nil, fmt.Errorf("invalid rule name '%s' or status", ctx.RuleName)
}

// createRule add a rule received from apiserver to this engine
func (p *ruleExecutor) createRule(ctx *RuleContext) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, found := p.rules[ctx.RuleName]; !found {
		if rule, err := newRule(p.config, ctx); err == nil {
			p.rules[ctx.RuleName] = rule
			return nil
		}
	}
	return fmt.Errorf("invalid rule name '%s' or status", ctx.RuleName)
}

// removeRule remove a rule from current rule engine
func (p *ruleExecutor) removeRule(r *RuleContext) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, found := p.rules[r.RuleName]; found {
		delete(p.rules, r.RuleName)
		return nil
	}
	return fmt.Errorf("rule '%s' no exist", r.RuleName)
}

// updateRule update rule in engine
func (p *ruleExecutor) updateRule(ctx *RuleContext) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if rule, err := p.getRule(ctx); err == nil {
		if _, found := p.rules[ctx.RuleName]; found {
			p.rules[ctx.RuleName] = rule
			return nil
		}
	}
	return fmt.Errorf("rule '%s' no exist", ctx.RuleName)
}

// startRule start rule in engine
func (p *ruleExecutor) startRule(r *RuleContext) error {
	// Check wether the rule executor is started
	if p.started == false {
		if err := p.start(); err != nil {
			return err
		}
		p.started = true
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()
	if rule, found := p.rules[r.RuleName]; found {
		rule.Status = ruleStatusStarted
		return nil
	}
	return fmt.Errorf("rule '%s' doesn't exist", r.RuleName)
}

// stopRule stop rule in engine
func (p *ruleExecutor) stopRule(r *RuleContext) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if rule, found := p.rules[r.RuleName]; found {
		rule.Status = ruleStatusStoped
		return nil
	}
	return fmt.Errorf("rule '%s' doesn't exist", r.RuleName)
}

// execute execute rule to process published topic data recevied from
// broker will be processed here and transformed into database
func (p *ruleExecutor) execute(e *event.TopicPublishEvent) error {
	glog.Info("executing rule ...")
	p.mutex.Lock()
	rules := p.rules
	p.mutex.Unlock()
	for _, rule := range rules {
		if rule.Status == ruleStatusStarted {
			if err := rule.handle(e); err != nil {
				glog.Infof("rule '%s' execution failed,'%s'", rule.RuleName, err.Error())
			}
		}
	}
	return nil
}
