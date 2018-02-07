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

package main

import (
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
)

var (
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
)

func initializeTestData(c config.Config) error {
	r, err := registry.New("conductor", c)
	if err != nil {
		return err
	}
	r.RegisterRule(&rule1)
	r.RegisterRule(&rule2)
	r.Close()
	return nil
}

func removeTestData(c config.Config) {
	if r, err := registry.New("conductor", c); err == nil {
		defer r.Close()
		r.RemoveRule(rule1.ProductId, rule1.RuleName)
		r.RemoveRule(rule2.ProductId, rule2.RuleName)
	}
}
