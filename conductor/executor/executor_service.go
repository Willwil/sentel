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
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/conductor/engine"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"
)

// executorService manage all rule engine and execute rule
type executorService struct {
	config       config.Config
	waitgroup    sync.WaitGroup
	quitChan     chan interface{}
	ruleChan     chan engine.RuleContext
	engines      map[string]*engine.RuleEngine
	mutex        sync.Mutex
	consumer     sarama.Consumer
	recoveryChan chan interface{}
}

const SERVICE_NAME = "executor"

type ServiceFactory struct{}

// New create executor service factory
func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	return &executorService{
		config:       c,
		waitgroup:    sync.WaitGroup{},
		quitChan:     make(chan interface{}),
		ruleChan:     make(chan engine.RuleContext),
		engines:      make(map[string]*engine.RuleEngine),
		mutex:        sync.Mutex{},
		recoveryChan: make(chan interface{}),
	}, nil
}

// Name
func (p *executorService) Name() string { return SERVICE_NAME }

// Initialize check all rule's state and recover if rule is started
func (p *executorService) Initialize() error {
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
func (p *executorService) Start() error {
	consumer, err := p.subscribeTopic(message.TopicNameRule)
	if err != nil {
		return fmt.Errorf("conductor failed to subscribe kafka event : %s", err.Error())
	}
	p.consumer = consumer
	// start rule channel
	p.waitgroup.Add(1)
	go func(s *executorService) {
		defer s.waitgroup.Done()
		select {
		case ctx := <-s.ruleChan:
			s.handleRule(ctx)
		case <-s.recoveryChan:
			s.recovery()
		case <-s.quitChan:
			break
		}
	}(p)
	return nil
}

// Stop
func (p *executorService) Stop() {
	p.quitChan <- true
	if p.consumer != nil {
		p.consumer.Close()
	}
	p.waitgroup.Wait()
	close(p.quitChan)
	close(p.ruleChan)
	close(p.recoveryChan)

	// stop all ruleEngine
	for _, engine := range p.engines {
		if engine != nil {
			engine.Stop()
		}
	}
}

// recovery restart rules after conductor is restarted
func (p *executorService) recovery() {
	r, _ := registry.New("conductor", p.config)
	defer r.Close()
	rules := r.GetRulesWithStatus(registry.RuleStatusStarted)
	for _, r := range rules {
		if r.Status == registry.RuleStatusStarted {
			ctx := engine.RuleContext{
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
func (p *executorService) handleRule(ctx engine.RuleContext) error {
	// Get engine instance according to product id
	if _, ok := p.engines[ctx.ProductId]; !ok { // not found
		ruleEngine, err := engine.NewRuleEngine(p.config, ctx.ProductId)
		if err != nil {
			glog.Errorf("Failed to create rule engint for product(%s)", ctx.ProductId)
			return err
		}
		p.engines[ctx.ProductId] = ruleEngine
	}
	r := p.engines[ctx.ProductId]

	switch ctx.Action {
	case engine.RuleCreate:
		return r.CreateRule(ctx)
	case engine.RuleRemove:
		return r.RemoveRule(ctx)
	case engine.RuleUpdate:
		return r.UpdateRule(ctx)
	case engine.RuleStart:
		return r.StartRule(ctx)
	case engine.RuleStop:
		return r.StopRule(ctx)
	}
	return nil
}

// subscribeTopc subscribe topics from apiserver
func (p *executorService) subscribeTopic(topic string) (sarama.Consumer, error) {
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

		go func(sarama.PartitionConsumer) {
			defer p.waitgroup.Done()
			for msg := range pc.Messages() {
				glog.Infof("donductor recevied topic '%s': '%s'", string(msg.Topic), msg.Value)
				p.handleNotifications(string(msg.Topic), msg.Value)
			}
		}(pc)
	}
	return consumer, nil
}

// handleNotifications handle notification from kafka
func (p *executorService) handleNotifications(topic string, value []byte) error {
	r := message.RuleTopic{}
	if err := json.Unmarshal(value, &r); err != nil {
		glog.Errorf("conductor failed to resolve topic from kafka: '%s'", err)
		return err
	}
	action := ""
	switch r.RuleAction {
	case message.RuleActionCreate:
		action = engine.RuleCreate
	case message.RuleActionRemove:
		action = engine.RuleRemove
	case message.RuleActionUpdate:
		action = engine.RuleUpdate
	case message.RuleActionStart:
		action = engine.RuleStart
	case message.RuleActionStop:
		action = engine.RuleStop
	default:
		return fmt.Errorf("Invalid rule action(%s) for product(%s)", r.RuleAction, r.ProductId)
	}
	ctx := engine.RuleContext{
		Action:    action,
		ProductId: r.ProductId,
		RuleName:  r.RuleName,
	}
	p.ruleChan <- ctx
	return nil
}

func HandleRuleNotification(ctx engine.RuleContext) {
	mgr := service.GetServiceManager()
	executor := mgr.GetService(SERVICE_NAME).(*executorService)
	executor.ruleChan <- ctx
}
