package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
)

const (
	fmtOfMqttEventBus   = "broker-%s-%s-event-mqtt"
	fmtOfBrokerEventBus = "broker-%s-%s-event-broker"
)

var (
	wg     sync.WaitGroup
	logger = log.New(os.Stderr, "[kafka]", log.LstdFlags)
)

func main() {
	tenant := flag.String("t", "ssddn", "Tenant")
	product := flag.String("p", "lock", "Product")
	flag.Parse()

	mqttBus := fmt.Sprintf(fmtOfMqttEventBus, *tenant, *product)
	brokerBus := fmt.Sprintf(fmtOfBrokerEventBus, *tenant, *product)

	sarama.Logger = logger

	consumer := make(map[string]sarama.Consumer)

	config := sarama.NewConfig()
	config.ClientID = "hello-mqtt"
	c1, err := sarama.NewConsumer(strings.Split("localhost:9092", ","), config)
	if err != nil {
		logger.Printf("Failed to start consumer: %s\n", err)
	}
	consumer[mqttBus] = c1

	config = sarama.NewConfig()
	config.ClientID = "hello-broker"
	c2, err := sarama.NewConsumer(strings.Split("localhost:9092", ","), config)
	if err != nil {
		logger.Printf("Failed to start consumer: %s\n", err)
	}
	consumer[brokerBus] = c2

	for topic, c := range consumer {
		partitionList, err := c.Partitions(topic)
		if err != nil {
			logger.Printf("Failed to get the list of partitions: %s\n", err)
		}

		for partition := range partitionList {
			pc, err := c.ConsumePartition(topic, int32(partition), -2)
			if err != nil {
				logger.Printf("Failed to start consumer for partition %d: %s\n", partition, err)
			}
			defer pc.AsyncClose()

			wg.Add(1)
			go func(topic string, pc sarama.PartitionConsumer) {
				defer wg.Done()
				for msg := range pc.Messages() {
					fmt.Printf("Topic:%s, Partition:%d, Offset:%d, Key:%s, Value:%s\n", topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				}
			}(topic, pc)
		}
	}

	wg.Wait()
	// consumer.Close()
}
