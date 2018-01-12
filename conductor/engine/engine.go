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
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/golang/glog"
)

// one product has one scpecific topic name
const fmtOfBrokerEventBus = "broker-%s-%s-event-broker"

// ruleEngine manage product's rules, add, start and stop rule
type RuleEngine struct {
	tenantId  string                // Tenant
	productId string                // one product have one rule engine
	rules     map[string]ruleWraper // all product's rule
	consumer  sarama.Consumer       // one product have one consumer
	config    config.Config         // configuration
	mutex     sync.Mutex            // mutex to protext rules list
	started   bool                  // indicate wether engined is started
	wg        sync.WaitGroup        // waitgroup for consumper partions
}

type RuleContext struct {
	Action    string `json:"action"`
	ProductId string `json:"productId"`
	RuleName  string `json:"ruleName"`
}

const (
	RuleCreate = "create"
	RuleRemove = "remove"
	RuleUpdate = "update"
	RuleStart  = "start"
	RuleStop   = "stop"
)
const (
	RuleStatusIdle    = "idle"
	RuleStatusStarted = "started"
	RuleStatusStoped  = "stoped"
)

// newRuleEngine create a engine according to product id and configuration
func NewRuleEngine(c config.Config, productId string) (*RuleEngine, error) {
	return &RuleEngine{
		productId: productId,
		config:    c,
		consumer:  nil,
		rules:     make(map[string]ruleWraper),
		mutex:     sync.Mutex{},
		wg:        sync.WaitGroup{},
		started:   false,
	}, nil
}

// start will start the rule engine, receiving topic and rule
func (p *RuleEngine) Start() error {
	if p.started {
		return fmt.Errorf("rule engine(%s) is already started", p.productId)
	}
	topic := fmt.Sprintf(fmtOfBrokerEventBus, p.tenantId, p.productId)
	if consumer, err := p.subscribeTopic(topic); err != nil {
		return err
	} else {
		p.consumer = consumer
		p.started = true
	}
	return nil
}

// subscribeTopc subscribe topics from apiserver
func (p *RuleEngine) subscribeTopic(topic string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.ClientID = "sentel_conductor_engine_" + p.productId
	khosts, _ := p.config.String("conductor", "kafka")
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("conductor failed to connect with kafka:%s", khosts)
	}
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		consumer.Close()
		return nil, fmt.Errorf("Failed to get list of partions:%s", err.Error())
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			consumer.Close()
			return nil, fmt.Errorf("Failed  to start consumer for partion %d:%s", partition, err)
		}
		p.wg.Add(1)

		go func(sarama.PartitionConsumer) {
			defer p.wg.Done()
			for msg := range pc.Messages() {
				obj := &event.MqttEvent{}
				if err := json.Unmarshal(msg.Value, obj); err == nil {
					p.handleKafkaEvent(obj)
				}
			}
		}(pc)
	}
	return consumer, nil
}

// handleEvent handle mqtt event from other service
func (p *RuleEngine) handleKafkaEvent(me *event.MqttEvent) error {
	e := &event.Event{}
	if err := json.Unmarshal(me.Common, &e.Common); err != nil {
		return fmt.Errorf("conductor unmarshal event common failed:%s", err.Error())
	}
	switch e.Common.Type {
	case event.SessionCreate:
	case event.SessionDestroy:
	case event.TopicSubscribe:
	case event.TopicUnsubscribe:
	case event.TopicPublish:
		detail := event.TopicPublishDetail{}
		if err := json.Unmarshal(me.Detail, &detail); err != nil {
			return fmt.Errorf("conductor unmarshal event detail failed:%s", err.Error())
		}
		e.Detail = detail
		return p.execute(e)
	}
	return nil
}

// stop will stop the engine
func (p *RuleEngine) Stop() {
	if p.consumer != nil {
		p.consumer.Close()
	}
	p.wg.Wait()
}

// getRuleObject get all rule's information from backend database
func (p *RuleEngine) getRuleObject(rc RuleContext) (*registry.Rule, error) {
	if r, err := registry.New("conductor", p.config); err == nil {
		defer r.Close()
		return r.GetRule(rc.RuleName, rc.ProductId)
	} else {
		return nil, err
	}
}

// createRule add a rule received from apiserver to this engine
func (p *RuleEngine) CreateRule(rc RuleContext) error {
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

// removeRule remove a rule from current rule engine
func (p *RuleEngine) RemoveRule(rc RuleContext) error {
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

// updateRule update rule in engine
func (p *RuleEngine) UpdateRule(rc RuleContext) error {
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

// startRule start rule in engine
func (p *RuleEngine) StartRule(rc RuleContext) error {
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
		p.rules[rc.RuleName].rule.Status = RuleStatusStarted
		return nil
	}
	return fmt.Errorf("rule:%s doesn't exist", rc.RuleName)
}

// stopRule stop rule in engine
func (p *RuleEngine) StopRule(rc RuleContext) error {
	glog.Infof("rule:%s is stoped", rc.RuleName)

	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.rules[rc.RuleName]; !ok { // not found
		return fmt.Errorf("Invalid rule:%s", rc.RuleName)
	}
	p.rules[rc.RuleName].rule.Status = RuleStatusStoped
	// Stop current engine if all rules are stoped
	for _, w := range p.rules {
		// If one of rule is not stoped, don't stop current engine
		if w.rule.Status != RuleStatusStoped {
			return nil
		}
	}
	p.Stop()
	return nil
}

// execute rule to process published topic
// Data recevied from iothub will be processed here and transformed into database
func (p *RuleEngine) execute(e *event.Event) error {
	for _, w := range p.rules {
		if w.rule.Status == RuleStatusStarted {
			if err := w.execute(p.config, e); err != nil {
				glog.Infof("conductor execute rule '%s' failed, reason: '%s'", w.rule.RuleName, err.Error())
			}
		}
	}
	return nil
}
