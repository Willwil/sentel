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
	"testing"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
)

var (
	defaultConfigs = config.M{
		"whaler": {
			"loglevel":        "debug",
			"kafka":           "localhost:9092",
			"mongo":           "localhost:27017",
			"connect_timeout": 5,
		},
	}

	rule1 = registry.Rule{
		ProductId:   "product1",
		RuleName:    "rule1",
		DataFormat:  "json",
		Description: "test rule1",
		DataProcess: registry.RuleDataProcess{
			Topic:     "hello",
			Condition: "",
			Fields:    []string{"field1", "field2"},
		},
		DataTarget: registry.RuleDataTarget{
			Type:  registry.DataTargetTypeTopic,
			Topic: "world",
		},
	}
	rule2 = registry.Rule{
		ProductId:   "product1",
		RuleName:    "rule2",
		DataFormat:  "json",
		Description: "test rule2",
		DataProcess: registry.RuleDataProcess{
			Topic:     "hello",
			Condition: "",
			Fields:    []string{"field1", "field2"},
		},
		DataTarget: registry.RuleDataTarget{
			Type:  registry.DataTargetTypeTopic,
			Topic: "world",
		},
	}

	defaultEngine *ruleEngine
	defaultConfig config.Config
)

func initializeTestData() error {
	r, err := registry.New(defaultConfig)
	if err != nil {
		return err
	}
	r.RegisterRule(&rule1)
	r.RegisterRule(&rule2)
	r.Close()
	return nil
}

func removeTestData() {
	if r, err := registry.New(defaultConfig); err == nil {
		defer r.Close()
		r.RemoveRule(rule1.ProductId, rule1.RuleName)
		r.RemoveRule(rule2.ProductId, rule2.RuleName)
	}
}

func Test_ruleEngine_Start(t *testing.T) {
	c := config.New("whaler")
	c.AddConfig(defaultConfigs)
	engine, err := newRuleEngine(c)
	if err != nil {
		t.Fatal(err)
	}
	defaultEngine = engine
	defaultConfig = c
	initializeTestData()
	if err := defaultEngine.Start(); err != nil {
		t.Fatal(err)
	}
}

func Test_ruleEngine_Recovery(t *testing.T) {
	if err := defaultEngine.Recovery(); err != nil {
		t.Error(err)
	}
}

func Test_ruleEngine_disptachRule(t *testing.T) {
	ctxs := []*RuleContext{
		NewRuleContext("product1", "rule1", RuleActionCreate),
		NewRuleContext("product1", "rule2", RuleActionCreate),
		NewRuleContext("product1", "rule1", RuleActionStart),
		NewRuleContext("product1", "rule2", RuleActionStart),
		NewRuleContext("product1", "rule1", RuleActionUpdate),
		NewRuleContext("product1", "rule2", RuleActionUpdate),
		NewRuleContext("product1", "rule1", RuleActionStop),
		NewRuleContext("product1", "rule2", RuleActionStop),
		NewRuleContext("product1", "rule1", RuleActionRemove),
		NewRuleContext("product1", "rule2", RuleActionRemove),
	}

	for _, ctx := range ctxs {
		if err := defaultEngine.dispatchRule(ctx); err != nil {
			t.Error(err)
		}
	}
}

func Test_ruleEngine_Stop(t *testing.T) {
	removeTestData()
	defaultEngine.Stop()
}
