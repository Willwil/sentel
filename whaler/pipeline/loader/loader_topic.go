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

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/whaler/pipeline/data"
)

type topicLoader struct {
	config   config.Config
	producer message.Producer
}

func newTopicLoader(c config.Config) (Loader, error) {
	producer, err := message.NewProducer(c, "", true)
	if err != nil {
		return nil, err
	}
	return &topicLoader{config: c, producer: producer}, nil
}

func (p *topicLoader) Name() string { return "topic" }
func (p *topicLoader) Close()       { p.producer.Close() }

func (p *topicLoader) Load(f *data.DataFrame) error {
	topic, ok1 := f.Context("topic").(string)
	e, ok2 := f.Context("event").(*event.TopicPublishEvent)
	if !ok1 || !ok2 || topic == "" || e == nil {
		return errors.New("invalid topic")
	}
	if buf, err := f.Serialize(); err != nil {
		return err
	} else {
		e.Payload = buf
		buf, _ := event.Encode(e, event.JSONCodec)
		msg := message.Broker{
			EventType: e.GetType(),
			Payload:   buf,
		}
		return p.producer.SendMessage(&msg)
	}
}
