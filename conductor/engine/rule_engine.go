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
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"
)

// ruleEngine manage all rule engine and execute rule
type ruleEngine struct {
	config       config.Config
	waitgroup    sync.WaitGroup
	quitChan     chan interface{}
	ruleChan     chan RuleContext
	executors    map[string]*ruleExecutor
	mutex        sync.Mutex
	consumer     sarama.Consumer
	recoveryChan chan interface{}
}

const SERVICE_NAME = "engine"

type ServiceFactory struct{}

// New create rule engine service factory
func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	return &ruleEngine{
		config:       c,
		waitgroup:    sync.WaitGroup{},
		quitChan:     make(chan interface{}),
		ruleChan:     make(chan RuleContext),
		executors:    make(map[string]*ruleExecutor),
		mutex:        sync.Mutex{},
		recoveryChan: make(chan interface{}),
	}, nil
}

// Name
func (p *ruleEngine) Name() string { return SERVICE_NAME }

// Initialize check all rule's state and recover if rule is started
func (p *ruleEngine) Initialize() error {
	// try connect with registry
	r, err := registry.New("conductor", p.config)
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
	consumer, err := p.subscribeTopic(message.TopicNameRule)
	if err != nil {
		return fmt.Errorf("conductor failed to subscribe kafka event : %s", err.Error())
	}
	p.consumer = consumer
	p.waitgroup.Add(1)
	go func(s *ruleEngine) {
		defer s.waitgroup.Done()
		for {
			select {
			case ctx := <-s.ruleChan:
				s.handleRule(ctx)
			case <-s.recoveryChan:
				s.recovery()
			case <-s.quitChan:
				return
			}
		}
	}(p)
	return nil
}

// Stop
func (p *ruleEngine) Stop() {
	p.quitChan <- true
	if p.consumer != nil {
		p.consumer.Close()
	}
	p.waitgroup.Wait()
	close(p.quitChan)
	close(p.ruleChan)
	close(p.recoveryChan)

	// stop all ruleEngine
	for _, executor := range p.executors {
		if executor != nil {
			executor.Stop()
		}
	}
}

// recovery restart rules after conductor is restarted
func (p *ruleEngine) recovery() {
	r, _ := registry.New("conductor", p.config)
	defer r.Close()
	rules := r.GetRulesWithStatus(registry.RuleStatusStarted)
	for _, r := range rules {
		if r.Status == registry.RuleStatusStarted {
			ctx := RuleContext{
				Action:    message.RuleActionStart,
				ProductId: r.ProductId,
				RuleName:  r.RuleName,
			}
			if err := p.handleRule(ctx); err != nil {
				glog.Errorf("product '%s', rule '%s'recovery failed", r.ProductId, r.RuleName)
			}
		}
	}
}

// handleRule process incomming rule
func (p *ruleEngine) handleRule(ctx RuleContext) error {
	// get executor instance according to product id
	if _, ok := p.executors[ctx.ProductId]; !ok { // not found
		executor, err := newRuleExecutor(p.config, ctx.ProductId)
		if err != nil {
			glog.Errorf("Failed to create rule engint for product(%s)", ctx.ProductId)
			return err
		}
		p.executors[ctx.ProductId] = executor
	}
	exe := p.executors[ctx.ProductId]

	switch ctx.Action {
	case RuleActionCreate:
		return exe.CreateRule(ctx)
	case RuleActionRemove:
		return exe.RemoveRule(ctx)
	case RuleActionUpdate:
		return exe.UpdateRule(ctx)
	case RuleActionStart:
		return exe.StartRule(ctx)
	case RuleActionStop:
		return exe.StopRule(ctx)
	}
	return nil
}

// subscribeTopc subscribe topics from apiserver
func (p *ruleEngine) subscribeTopic(topic string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.ClientID = "sentel_conductor_indicator"
	khosts, _ := p.config.String("conductor", "kafka")
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", khosts)
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
		p.waitgroup.Add(1)

		go func(p *ruleEngine, pc sarama.PartitionConsumer) {
			defer p.waitgroup.Done()
			for msg := range pc.Messages() {
				glog.Infof("donductor recevied topic '%s': '%s'", string(msg.Topic), msg.Value)
				p.handleNotifications(string(msg.Topic), msg.Value)
			}
		}(p, pc)
	}
	return consumer, nil
}

// handleNotifications handle notification from kafka
func (p *ruleEngine) handleNotifications(topic string, value []byte) error {
	r := message.RuleTopic{}
	if err := json.Unmarshal(value, &r); err != nil {
		glog.Errorf("conductor failed to resolve topic from kafka: '%s'", err)
		return err
	}
	action := ""
	switch r.RuleAction {
	case message.RuleActionCreate:
		action = RuleActionCreate
	case message.RuleActionRemove:
		action = RuleActionRemove
	case message.RuleActionUpdate:
		action = RuleActionUpdate
	case message.RuleActionStart:
		action = RuleActionStart
	case message.RuleActionStop:
		action = RuleActionStop
	default:
		return fmt.Errorf("Invalid rule action(%s) for product(%s)", r.RuleAction, r.ProductId)
	}
	ctx := RuleContext{
		Action:    action,
		ProductId: r.ProductId,
		RuleName:  r.RuleName,
	}
	p.ruleChan <- ctx
	return nil
}

func HandleRuleNotification(ctx RuleContext) {
	mgr := service.GetServiceManager()
	s := mgr.GetService(SERVICE_NAME).(*ruleEngine)
	s.ruleChan <- ctx
}
