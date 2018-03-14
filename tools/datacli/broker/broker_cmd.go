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

package broker

import (
	"fmt"

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var (
	tenantId  string
	topicName string
	payload   string
)

var command = &cobra.Command{
	Use:   "broker",
	Short: "send test data to kafka as broker.",
	Run:   handleCommand,
}

func handleCommand(cmd *cobra.Command, args []string) {
	c := config.New("datacli")
	c.AddConfigItem("kafka", "localhost:9092")
	producer, err := message.NewProducer(c, "datacli", true)
	if err != nil {
		glog.Info(err)
		return
	}
	defer producer.Close()
	e := event.TopicPublishEvent{
		Type:    event.TopicPublish,
		Topic:   topicName,
		Payload: []byte(payload),
		Qos:     1,
		Retain:  true,
	}
	topic := fmt.Sprintf(event.FmtOfBrokerEventBus, tenantId)
	value, _ := event.Encode(&e, event.JSONCodec)
	msg := message.Broker{TopicName: topic, Payload: value}
	if err := producer.SendMessage(&msg); err != nil {
		glog.Info(err)
	}

}

func NewCommand() *cobra.Command {
	command.Flags().StringVarP(&tenantId, "tenant", "t", "", "tenant")
	command.Flags().StringVarP(&topicName, "topic", "n", "", "topic name")
	command.Flags().StringVarP(&payload, "payload", "p", "", "payload")

	return command
}
