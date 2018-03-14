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

package extractor

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/cloustone/sentel/whaler/pipeline/data"
	"github.com/elgs/jsonql"
)

type eventExtractor struct {
	dataprocess registry.RuleDataProcess
	config      config.Config
}

func newEventExtractor(c config.Config) (Extractor, error) {
	return &eventExtractor{
		config:      c,
		dataprocess: c.MustValue("dataprocess").(registry.RuleDataProcess),
	}, nil
}

func (p *eventExtractor) Extract(origin interface{}) (*data.DataFrame, error) {
	if _, ok := origin.(*event.TopicPublishEvent); !ok {
		return nil, errors.New("invalid parameter")
	}

	e := origin.(*event.TopicPublishEvent)
	f := data.NewDataFrame()
	if e.Topic == p.dataprocess.Topic {
		f.SetContext("topic", e.Topic)
		f.SetContext("event", e)
		if p.dataprocess.Condition == "" {
			// return data if no process condition is specified
			err := json.Unmarshal(e.Payload, &f)
			return f, err
		} else {
			// topic's data must be json format
			parser, err := jsonql.NewStringQuery(string(e.Payload))
			if err != nil {
				return nil, fmt.Errorf("data format is invalid '%s'", err.Error())
			}
			// parse condiction and get result
			n, err := parser.Query(p.dataprocess.Condition)
			if err != nil || n == nil {
				return nil, fmt.Errorf("data query failed '%s'", err.Error())
			}
			switch n.(type) {
			case map[string]interface{}:
				results := n.(map[string]interface{})
				fields := p.dataprocess.Fields
				for _, field := range fields {
					if _, found := results[field]; found { // message field
						f.AddField(field, results[field])
					} else if len(field) > 0 && field[0] == '$' { // variable support
						v := p.getVariableValue(field[1:], e)
						f.AddField(field, v)
					} else { // functions support
						if n, err := p.getFunctionValue(field, e); err == nil {
							f.AddField(field, n)
						}
					}
				}
				return f, nil
			}
		}
	}
	return nil, errors.New("data format is invalid to query")
}

func (p *eventExtractor) getFunctionValue(name string, e *event.TopicPublishEvent) (interface{}, error) {
	return nil, errors.New("not implemented")
}

func (p *eventExtractor) getVariableValue(name string, e *event.TopicPublishEvent) interface{} {
	switch name {
	case "deviceId":
		return e.ClientID
	}
	return nil
}

func (p *eventExtractor) Close() {}
