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

// one product has one scpecific topic name
const fmtOfBrokerEventBus = "iot-%s-%s-event-broker"

// ruleEngine manage product's rules, add, start and stop rule
type ruleExecutor struct {
	config    config.Config          // configuration
	tenantId  string                 // Tenant
	productId string                 // one product have one rule engine
	rules     map[string]*ruleWraper // all product's rule
	mutex     sync.Mutex             // mutex to protext rules list
	started   bool                   // indicate wether engined is started
	consumer  message.Consumer
	dataChan  chan event.Event
	quitChan  chan interface{}
	waitgroup sync.WaitGroup
}

type RuleContext struct {
	Action    string `json:"action"`
	ProductId string `json:"productId"`
	RuleName  string `json:"ruleName"`
	Resp      chan error
}

const (
	RuleActionCreate = "create"
	RuleActionRemove = "remove"
	RuleActionUpdate = "update"
	RuleActionStart  = "start"
	RuleActionStop   = "stop"
)
const (
	ruleStatusIdle    = "idle"
	ruleStatusStarted = "started"
	ruleStatusStoped  = "stoped"
)

// newRuleExecutor create a engine according to product id and configuration
func newRuleExecutor(c config.Config, productId string) (*ruleExecutor, error) {
	khosts := c.MustString("conductor", "kafka")
	consumer, _ := message.NewConsumer(khosts, "conductorRuleExecutor")
	return &ruleExecutor{
		productId: productId,
		config:    c,
		rules:     make(map[string]*ruleWraper),
		mutex:     sync.Mutex{},
		started:   false,
		consumer:  consumer,
		dataChan:  make(chan event.Event, 1),
		quitChan:  make(chan interface{}, 1),
		waitgroup: sync.WaitGroup{},
	}, nil
}

// start will start the rule engine, receiving topic and rule
func (p *ruleExecutor) start() error {
	if p.started {
		return fmt.Errorf("rule executor '%s' already started", p.productId)
	}
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

	topic := fmt.Sprintf(fmtOfBrokerEventBus, p.tenantId, p.productId)
	if err := p.consumer.Subscribe(topic, p.messageHandlerFunc, nil); err != nil {
		return err
	}
	p.consumer.Start()
	p.started = true
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
func (p *ruleExecutor) messageHandlerFunc(topic string, value []byte, ctx interface{}) {
	t, err := event.Decode(value, nil)
	if err == nil && t != nil && t.GetType() == event.TopicPublish {
		p.dataChan <- t
		// we can call p.execute(t) here, but in consider of avoiding to block message receiver
		// we use datachannel
	}
}

// getRuleObject get all rule's information from backend database
func (p *ruleExecutor) getRuleObject(rc RuleContext) (*registry.Rule, error) {
	if r, err := registry.New("conductor", p.config); err == nil {
		defer r.Close()
		return r.GetRule(rc.ProductId, rc.RuleName)
	} else {
		return nil, err
	}
}

// createRule add a rule received from apiserver to this engine
func (p *ruleExecutor) createRule(r RuleContext) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, found := p.rules[r.RuleName]; !found {
		if obj, err := p.getRuleObject(r); err == nil {
			if w, err := newRuleWraper(p.config, obj); err == nil {
				p.rules[r.RuleName] = w
				glog.Infof("rule '%s' is added", r.RuleName)
				return nil
			}
			return fmt.Errorf("rule wrappr '%s' failure", r.RuleName)
		}
		return fmt.Errorf("rule '%s' metadata fetch failure", r.RuleName)
	}
	return fmt.Errorf("rule '%s' already exist", r.RuleName)
}

// removeRule remove a rule from current rule engine
func (p *ruleExecutor) removeRule(r RuleContext) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, found := p.rules[r.RuleName]; found {
		delete(p.rules, r.RuleName)
		glog.Infof("rule '%s' is removed", r.RuleName)
		return nil
	}
	return fmt.Errorf("rule '%s' no exist", r.RuleName)
}

// updateRule update rule in engine
func (p *ruleExecutor) updateRule(r RuleContext) error {
	if obj, err := p.getRuleObject(r); err == nil {
		p.mutex.Lock()
		defer p.mutex.Unlock()
		if w, found := p.rules[r.RuleName]; found {
			w.rule = obj
			p.rules[r.RuleName] = w
			glog.Infof("rule '%s' is updated", r.RuleName)
			return nil
		}
	}
	return fmt.Errorf("rule '%s' no exist", r.RuleName)
}

// startRule start rule in engine
func (p *ruleExecutor) startRule(r RuleContext) error {
	// Check wether the rule executor is started
	if p.started == false {
		if err := p.start(); err != nil {
			glog.Errorf("start rule '%s' failed,'%s'", r.RuleName, err.Error())
			return err
		}
		p.started = true
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()
	if w, found := p.rules[r.RuleName]; found {
		w.rule.Status = ruleStatusStarted
		glog.Infof("rule '%s' is started", r.RuleName)
		return nil
	}
	return fmt.Errorf("rule '%s' doesn't exist", r.RuleName)
}

// stopRule stop rule in engine
func (p *ruleExecutor) stopRule(r RuleContext) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if w, found := p.rules[r.RuleName]; !found { // not found
		w.rule.Status = ruleStatusStoped
		return nil
	}
	return fmt.Errorf("rule '%s' doesn't exist", r.RuleName)
}

// execute rule to process published topic
// Data recevied from iothub will be processed here and transformed into database
func (p *ruleExecutor) execute(e event.Event) error {
	p.mutex.Lock()
	rules := p.rules
	p.mutex.Unlock()
	for _, w := range rules {
		if w.rule.Status == ruleStatusStarted {
			if err := w.handle(e); err != nil {
				glog.Infof("rule '%s' execution failed,'%s'", w.rule.RuleName, err.Error())
			}
		}
	}
	return nil
}
