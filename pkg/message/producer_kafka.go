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
package message

import (
	"errors"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/golang/glog"
)

type kafkaProducer struct {
	khosts        string
	clientID      string
	sync          bool
	syncProducer  sarama.SyncProducer
	asyncProducer sarama.AsyncProducer
}

func newKafkaProducer(c config.Config, clientID string, sync bool) (Producer, error) {
	khosts, err := c.String("kafka")
	if err != nil || khosts == "" {
		return nil, errors.New("invalid kafka setting")
	}
	names := strings.Split(khosts, ":")
	if len(names) == 1 {
		khosts = khosts + ":9092"
	}

	p := &kafkaProducer{
		khosts:   khosts,
		clientID: clientID,
		sync:     sync,
	}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.ClientID = "sentel"

	if sync {
		producer, err := sarama.NewSyncProducer(strings.Split(khosts, ","), config)
		if err != nil {
			glog.Error(err)
			return nil, err
		}
		p.syncProducer = producer
	} else {
		producer, err := sarama.NewAsyncProducer(strings.Split(khosts, ","), config)
		if err != nil {
			glog.Error(err)
			return nil, err
		}
		p.asyncProducer = producer
	}
	return p, nil
}

func (p *kafkaProducer) SendMessage(t Message) error {
	value, err := t.Serialize(JSONSerialization)
	if err != nil {
		return err
	}
	topic := t.Topic()
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}

	if p.sync && p.syncProducer != nil {
		if _, _, err := p.syncProducer.SendMessage(msg); err != nil {
			return err
		}
	} else if p.asyncProducer != nil {
		go func(p sarama.AsyncProducer) {
			errors := p.Errors()
			success := p.Successes()
			select {
			case err := <-errors:
				if err != nil {
					glog.Error(err)
				}
			case <-success:
			}
		}(p.asyncProducer)
		p.asyncProducer.Input() <- msg
	} else {
		return errors.New("invalid producer")
	}
	return nil
}

func (p *kafkaProducer) SendMessages(msgs []Message) error {
	kmsgs := []*sarama.ProducerMessage{}
	for _, msg := range msgs {
		value, err := msg.Serialize(JSONSerialization)
		if err != nil {
			return err
		}
		topic := msg.Topic()
		kmsg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(value),
		}
		kmsgs = append(kmsgs, kmsg)
	}

	if p.sync && p.syncProducer != nil {
		return p.syncProducer.SendMessages(kmsgs)
	} else if p.asyncProducer != nil {
		go func(p sarama.AsyncProducer) {
			errors := p.Errors()
			success := p.Successes()
			select {
			case err := <-errors:
				if err != nil {
					glog.Error(err)
				}
			case <-success:
			}
		}(p.asyncProducer)
		for _, msg := range kmsgs {
			p.asyncProducer.Input() <- msg
		}
	} else {
		return errors.New("invalid producer")
	}
	return nil
}

func (p *kafkaProducer) Close() {
	if p.asyncProducer != nil {
		p.asyncProducer.Close()
	} else if p.syncProducer != nil {
		p.syncProducer.Close()
	}
}
