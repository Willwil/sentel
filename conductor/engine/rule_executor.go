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
	"encoding/json"
	"fmt"
	"sync"

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/golang/glog"
)

// one product has one scpecific topic name
const fmtOfBrokerEventBus = "broker-%s-%s-event-broker"

// ruleEngine manage product's rules, add, start and stop rule
type ruleExecutor struct {
	tenantId  string                // Tenant
	productId string                // one product have one rule engine
	rules     map[string]ruleWraper // all product's rule
	config    config.Config         // configuration
	mutex     sync.Mutex            // mutex to protext rules list
	started   bool                  // indicate wether engined is started
	consumer  message.Consumer
}

type RuleContext struct {
	Action    string
	ProductId string `json:"productId"`
	RuleName  string `json:"ruleName"`
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
		rules:     make(map[string]ruleWraper),
		mutex:     sync.Mutex{},
		started:   false,
		consumer:  consumer,
	}, nil
}

// Start will start the rule engine, receiving topic and rule
func (p *ruleExecutor) Start() error {
	if p.started {
		return fmt.Errorf("rule engine(%s) is already started", p.productId)
	}
	topic := fmt.Sprintf(fmtOfBrokerEventBus, p.tenantId, p.productId)
	if err := p.consumer.Subscribe(topic, p.messageHandlerFunc, nil); err != nil {
		return err
	}
	p.consumer.Start()
	p.started = true
	return nil
}

// Stop will stop the engine
func (p *ruleExecutor) Stop() {
	p.consumer.Close()
}

// messageHandlerFunc handle mqtt event from other service
func (p *ruleExecutor) messageHandlerFunc(topic string, value []byte, ctx interface{}) {
	re := event.RawEvent{}
	if err := json.Unmarshal(value, &re); err != nil {
		glog.Errorf("conductor unmarshal event common failed:%s", err.Error())
		return
	}

	e := &event.Event{}
	if err := json.Unmarshal(re.Header, &e.EventHeader); err != nil {
		glog.Errorf("conductor unmarshal event common failed:%s", err.Error())
		return
	}
	switch e.Type {
	case event.SessionCreate:
	case event.SessionDestroy:
	case event.TopicSubscribe:
	case event.TopicUnsubscribe:
	case event.TopicPublish:
		detail := event.TopicPublishDetail{}
		if err := json.Unmarshal(re.Payload, &detail); err != nil {
			glog.Errorf("conductor unmarshal event detail failed:%s", err.Error())
		}
		e.Detail = detail
		p.execute(e)
	}
}

// getRuleObject get all rule's information from backend database
func (p *ruleExecutor) getRuleObject(rc RuleContext) (*registry.Rule, error) {
	if r, err := registry.New("conductor", p.config); err == nil {
		defer r.Close()
		return r.GetRule(rc.RuleName, rc.ProductId)
	} else {
		return nil, err
	}
}

// CreateRule add a rule received from apiserver to this engine
func (p *ruleExecutor) CreateRule(rc RuleContext) error {
	glog.Infof("rule:%s is added", rc.RuleName)
	obj, err := p.getRuleObject(rc)
	if err != nil {
		return err
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.rules[rc.RuleName]; ok {
		return fmt.Errorf("rule:%s already exist", rc.RuleName)
	}
	p.rules[rc.RuleName] = ruleWraper{rule: obj}
	return nil
}

// RemoveRule remove a rule from current rule engine
func (p *ruleExecutor) RemoveRule(rc RuleContext) error {
	glog.Infof("Rule:%s is deleted", rc.RuleName)

	// Get rule detail
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.rules[rc.RuleName]; ok {
		delete(p.rules, rc.RuleName)
		return nil
	}
	return fmt.Errorf("rule:%s doesn't exist", rc.RuleName)
}

// UpdateRule update rule in engine
func (p *ruleExecutor) UpdateRule(rc RuleContext) error {
	glog.Infof("Rule:%s is updated", rc.RuleName)

	obj, err := p.getRuleObject(rc)
	if err != nil {
		return err
	}

	// Get rule detail
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.rules[rc.RuleName]; ok {
		p.rules[rc.RuleName] = ruleWraper{rule: obj}
		return nil
	}
	return fmt.Errorf("rule:%s doesn't exist", rc.RuleName)
}

// StartRule start rule in engine
func (p *ruleExecutor) StartRule(rc RuleContext) error {
	glog.Infof("rule:%s is started", rc.RuleName)

	// Check wether the rule engine is started
	if p.started == false {
		if err := p.Start(); err != nil {
			glog.Errorf(" conductor start rule '%s' failed:'%v'", rc.RuleName, err)
			return err
		}
		p.started = true
	}

	// Start the rule
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.rules[rc.RuleName]; ok {
		p.rules[rc.RuleName].rule.Status = ruleStatusStarted
		return nil
	}
	return fmt.Errorf("rule:%s doesn't exist", rc.RuleName)
}

// StopRule stop rule in engine
func (p *ruleExecutor) StopRule(rc RuleContext) error {
	glog.Infof("rule:%s is stoped", rc.RuleName)

	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.rules[rc.RuleName]; !ok { // not found
		return fmt.Errorf("Invalid rule:%s", rc.RuleName)
	}
	p.rules[rc.RuleName].rule.Status = ruleStatusStoped
	// Stop current engine if all rules are stoped
	for _, w := range p.rules {
		// If one of rule is not stoped, don't stop current engine
		if w.rule.Status != ruleStatusStoped {
			return nil
		}
	}
	p.Stop()
	return nil
}

// execute rule to process published topic
// Data recevied from iothub will be processed here and transformed into database
func (p *ruleExecutor) execute(e *event.Event) error {
	for _, w := range p.rules {
		if w.rule.Status == ruleStatusStarted {
			if err := w.execute(p.config, e); err != nil {
				glog.Infof("conductor execute rule '%s' failed, reason: '%s'", w.rule.RuleName, err.Error())
			}
		}
	}
	return nil
}
