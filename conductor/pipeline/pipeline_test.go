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

package pipeline

import (
	"testing"

	"github.com/cloustone/sentel/broker/event"
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
)

type mockReader struct{}

func (p *mockReader) Read() (interface{}, error) {
	return &event.Event{
		EventHeader: event.EventHeader{Type: event.TopicPublish, ClientId: "jenson"},
		Detail: &event.TopicPublishDetail{
			Topic:   "hello",
			Payload: []byte("hello, world"),
		},
	}, nil
}

func TestETL_Run(t *testing.T) {
	config := config.New()
	options := make(map[string]string)
	options["extractor"] = "event"
	options["transformer"] = "no"
	options["loader"] = "topic"
	options["mongo"] = "localhost:27017"
	options["kafka"] = "localhost:9092"
	config.AddConfigSection("etl", options)
	etl, err := New(config)
	if err != nil {
		t.Error(err)
		return
	}
	if err := etl.Run(&mockReader{}, &rule1); err != nil {
		t.Error(err)
	}
}
