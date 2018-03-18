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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/whaler/pipeline/data"
	"github.com/golang/glog"
)

type topicLoader struct {
	config   config.Config
	producer message.Producer
	tenantId string
}

func newTopicLoader(c config.Config) (Loader, error) {
	producer, err := message.NewProducer(c, "", true)
	if err != nil {
		return nil, err
	}
	return &topicLoader{
		config:   c,
		producer: producer,
		tenantId: c.MustString("tenantId"),
	}, nil
}

func (p *topicLoader) Name() string { return "topic" }
func (p *topicLoader) Close()       { p.producer.Close() }

func (p *topicLoader) Load(f *data.DataFrame) error {
	if e, _ := f.Context("event").(*event.TopicPublishEvent); e != nil {
		if buf, err := f.Serialize(); err == nil {
			dstTopic := p.config.MustString("topic")
			glog.Infof("data from '%s' topic is transfered to '%s' topic", e.Topic, dstTopic)
			e.Payload = buf
			e.Topic = dstTopic
			topic := fmt.Sprintf(event.FmtOfBrokerEventBus, p.tenantId)
			value, _ := json.Marshal(&e)
			msg := message.Broker{EventType: event.TopicPublish, TopicName: topic, Payload: value}
			return p.producer.SendMessage(&msg)
		}
	}
	return errors.New("invalid data frame")
}
