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

package loader

import (
	"errors"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
)

type topicLoader struct {
	config   config.Config
	producer message.Producer
}

func newTopicLoader(c config.Config) (Loader, error) {
	hosts, err := c.String("etl", "message_server")
	if err != nil || hosts == "" {
		return nil, errors.New("invalid message server setting")
	}
	producer, err := message.NewProducer(hosts, "", true)
	if err != nil {
		return nil, err
	}
	return &topicLoader{config: c, producer: producer}, nil
}

func (p *topicLoader) Name() string { return "topic" }
func (p *topicLoader) Close() {
	p.producer.Close()
}

func (p *topicLoader) Load(data map[string]interface{}, ctx interface{}) error {
	topic, err := p.config.String("etl", "topic")
	if err != nil || topic == "" {
		return errors.New("invalid topic")
	}
	return p.producer.SendMessage("", topic, data)
}
