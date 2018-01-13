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

package data

import (
	"errors"
	"fmt"

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/elgs/jsonql"
)

type DataProcessor interface {
	Execute(e *event.Event) (map[string]interface{}, error)
}

func NewProcessor(c config.Config, r *registry.Rule) (DataProcessor, error) {
	return &simpleDataProcessor{config: c, rule: r}, nil
}

type simpleDataProcessor struct {
	config config.Config
	rule   *registry.Rule
}

func (p *simpleDataProcessor) Execute(e *event.Event) (map[string]interface{}, error) {
	detail := e.Detail.(*event.TopicPublishDetail)
	if detail.Topic == p.rule.DataProcess.Topic {
		// topic's data must be json format
		parser, err := jsonql.NewStringQuery(string(detail.Payload))
		if err != nil {
			return nil, fmt.Errorf("data format is invalid '%s'", err.Error())
		}
		// parse condiction and get result
		n, err := parser.Query(p.rule.DataProcess.Condition)
		if err != nil {
			return nil, fmt.Errorf("data query failed '%s'", err.Error())
		}
		switch n.(type) {
		case map[string]interface{}:
			results := n.(map[string]interface{})
			fields := p.rule.DataProcess.Fields
			data := make(map[string]interface{})
			for _, field := range fields {
				if _, found := results[field]; found { // message field
					data[field] = results[field]
				} else if len(field) > 0 && field[0] == '$' { // variable support
					n := field[1:]
					data[n] = p.getVariableValue(n, e)
				} else { // functions support
					if n, err := p.getFunctionValue(field, e); err == nil {
						data[field] = n
					}
				}
			}
			return data, nil
		}
	}
	return nil, errors.New("data format is invalid to query")
}

func (p *simpleDataProcessor) getFunctionValue(name string, e *event.Event) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (p *simpleDataProcessor) getVariableValue(name string, e *event.Event) interface{} {
	switch name {
	case "deviceId":
		return e.ClientId
	}
	return nil
}
