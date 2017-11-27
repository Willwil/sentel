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
package util

import (
//	"errors"
	"strings"
	"time"
	"fmt"

	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

var kafka string
/*
run kafka.
step 1:(must)
bin/zookeeper-server-start.sh config/zookeeper.properties

step 2:(must)
bin/kafka-server-start.sh config/server.properties

step 3:create topic(must),ie:product
kafka-topics.sh --create --replication-factor 1 --partitions 8 --topic product --zookeeper localhost:2181

step 4:list.
kafka-topics.sh --list --zookeeper localhost:2181

step 5:send msg
kafka-console-producer.sh --broker-list localhost:9092 --topic product

step 6:recv msg
kafka-console-consumer.sh --topic product --zookeeper localhost:2181 --from-beginning


test:
run sentel-apiserver
run sentel-conductor
curl -d "Description=7&Name=2" "http://localhost:4145/api/v1/products/3?api-version=v1"
conductor output msg.
*/

func InitializeKafka(c core.Config) {
	kafka = c.MustString("apiserver", "kafka")
	fmt.Printf("Initializing kafka:%s...", kafka)
}

func SyncProduceMessage(cfg core.Config, topic string, key string, value sarama.Encoder) error {
	// Get kafka server
	glog.Infof("SyncProduceMessage cfg:[%s] :%s, %s, %s ", cfg, topic, key, value)
//	kafka, err := cfg.String("kafka", "hosts")
//	if err != nil || kafka == "" {
//		return errors.New("Invalid kafka configuration")
//	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: value,
	}

	producer, err := sarama.NewSyncProducer(strings.Split(kafka, ","), config)
	if err != nil {
		glog.Errorf("Failed to produce message:%s", err.Error())
		return err
	}
	defer producer.Close()

	if _, _, err := producer.SendMessage(msg); err != nil {
		glog.Errorf("Failed to send producer message:%s", err.Error())
	}
	return err
}

func AsyncProduceMessage(cfg core.Config, key string, topic string, value sarama.Encoder) error {
	// Get kafka server
	fmt.Printf("AsyncProduceMessage1 cfg:[%s] :%s, %s, %s \n", cfg, topic, key, value)
	glog.Infof("AsyncProduceMessage cfg:[%s] :%s, %s ", cfg, topic, key)
	//kafka, err := cfg.String("kafka", "hosts")
	//if err != nil || kafka == "" {
	//fmt.Printf("error cfg:[%s] :%s \n", kafka, err)
	//	return errors.New("Invalid kafka configuration")
	//}
	fmt.Printf(" cfg:[%s] \n", kafka)

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second

	//v := "test kafka"
	v,_ := json.Marshal(value)
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(v),//value,
	}

	producer, err := sarama.NewAsyncProducer(strings.Split(kafka, ","), config)
	if err != nil {
		glog.Errorf("Failed to produce message:%s", err.Error())
		return err
	}
	defer producer.Close()

	go func(p sarama.AsyncProducer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					glog.Error(err)
				}
			case <-success:
			}
		}
	}(producer)

	producer.Input() <- msg
	return err
}
