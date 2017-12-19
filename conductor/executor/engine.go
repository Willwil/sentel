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

package executor

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/core"
	"github.com/cloustone/sentel/core/db"
	"github.com/golang/glog"
)

// one product has one scpecific topic name
const publishTopicUrl = "/cluster/publish/%s"

// publishTopic is topic object received from kafka
type publishTopic struct {
	ClientId  string `json:"clientId"`
	Topic     string `json:"topic"`
	ProductId string `json:"product"`
	Content   string `json:"content"`
	encoded   []byte
	err       error
}

// ruleEngine manage product's rules, add, start and stop rule
type ruleEngine struct {
	productId string              // one product have one rule engine
	rules     map[string]*db.Rule // all product's rule
	consumer  sarama.Consumer     // one product have one consumer
	config    core.Config         // configuration
	mutex     sync.Mutex          // mutex to protext rules list
	started   bool                // indicate wether engined is started
	wg        sync.WaitGroup      // waitgroup for consumper partions
}

// newRuleEngine create a engine according to product id and configuration
func newRuleEngine(c core.Config, productId string) (*ruleEngine, error) {
	return &ruleEngine{
		productId: productId,
		config:    c,
		consumer:  nil,
		rules:     make(map[string]*db.Rule),
		mutex:     sync.Mutex{},
		wg:        sync.WaitGroup{},
		started:   false,
	}, nil
}

// start will start the rule engine, receiving topic and rule
func (p *ruleEngine) start() error {
	if p.started {
		return fmt.Errorf("rule engine(%s) is already started", p.productId)
	}
	topic := fmt.Sprintf(publishTopicUrl, p.productId)
	if consumer, err := p.subscribeTopic(topic); err != nil {
		return err
	} else {
		p.consumer = consumer
		p.started = true
	}
	return nil
}

// subscribeTopc subscribe topics from apiserver
func (p *ruleEngine) subscribeTopic(topic string) (sarama.Consumer, error) {
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
				t := publishTopic{}
				if err := json.Unmarshal(msg.Value, &t); err != nil {
					glog.Errorf("Failed to handle topic:%v", err)
					continue
				}
				if err := p.execute(&t); err != nil {
					glog.Errorf("Failed to handle topic:%v", err)
					continue
				}
			}
		}(pc)
	}
	return consumer, nil
}

// stop will stop the engine
func (p *ruleEngine) stop() {
	if p.consumer != nil {
		p.consumer.Close()
	}
	p.wg.Wait()
}

// getRuleObject get all rule's information from backend database
func (p *ruleEngine) getRuleObject(productId, ruleName string) (*db.Rule, error) {
	if registry, err := db.NewRegistry("conductor", p.config); err == nil {
		defer registry.Release()
		return registry.GetRule(ruleName, productId)
	} else {
		return nil, err
	}
}

// createRule add a rule received from apiserver to this engine
func (p *ruleEngine) createRule(r *db.Rule) error {
	glog.Infof("ruld:%s is added", r.RuleName)
	obj, err := p.getRuleObject(r.ProductId, r.RuleName)
	if err != nil {
		return err
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.rules[r.RuleName]; ok {
		return fmt.Errorf("rule:%s already exist", r.RuleName)
	}
	p.rules[r.RuleName] = obj
	return nil
}

// removeRule remove a rule from current rule engine
func (p *ruleEngine) removeRule(r *db.Rule) error {
	glog.Infof("Rule:%s is deleted", r.RuleName)

	// Get rule detail
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.rules[r.RuleName]; ok {
		delete(p.rules, r.RuleName)
		return nil
	}
	return fmt.Errorf("rule:%s doesn't exist", r.RuleName)
}

// updateRule update rule in engine
func (p *ruleEngine) updateRule(r *db.Rule) error {
	glog.Infof("Rule:%s is updated", r.RuleName)

	obj, err := p.getRuleObject(r.ProductId, r.RuleName)
	if err != nil {
		return err
	}

	// Get rule detail
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.rules[r.RuleName]; ok {
		p.rules[r.RuleName] = obj
		return nil
	}
	return fmt.Errorf("rule:%s doesn't exist", r.RuleName)
}

// startRule start rule in engine
func (p *ruleEngine) startRule(r *db.Rule) error {
	glog.Infof("rule:%s is started", r.RuleName)

	// Check wether the rule engine is started
	if p.started == false {
		if err := p.start(); err != nil {
			glog.Errorf(" conductor start rule '%s' failed:'%v'", r.RuleName, err)
			return err
		}
		p.started = true
	}

	// Start the rule
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.rules[r.RuleName]; ok {
		p.rules[r.RuleName].Status = RuleStatusStarted
		return nil
	}
	return fmt.Errorf("rule:%s doesn't exist", r.RuleName)
}

// stopRule stop rule in engine
func (p *ruleEngine) stopRule(r *db.Rule) error {
	glog.Infof("rule:%s is stoped", r.RuleName)

	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.rules[r.RuleName]; !ok { // not found
		return fmt.Errorf("Invalid rule:%s", r.RuleName)
	}
	p.rules[r.RuleName].Status = RuleStatusStoped
	// Stop current engine if all rules are stoped
	for _, rule := range p.rules {
		// If one of rule is not stoped, don't stop current engine
		if rule.Status != RuleStatusStoped {
			return nil
		}
	}
	p.stop()
	return nil
}

// execute rule to process published topic
// Data recevied from iothub will be processed here and transformed into database
func (p *ruleEngine) execute(t *publishTopic) error {
	return nil
}
