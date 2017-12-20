package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

const (
	fmtOfMqttEventBus   = "broker-%s-%s-event-mqtt"
	fmtOfBrokerEventBus = "broker-%s-%s-event-broker"
)

var (
	logger   = log.New(os.Stderr, "[kafka]", log.LstdFlags)
	producer sarama.AsyncProducer
)

func newProducer() sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	p, err := sarama.NewAsyncProducer(strings.Split("localhost:9092", ","), config)
	if err != nil {
		logger.Printf("Failed to create producer: %s\n", err)
		os.Exit(500)
	}

	go func() {
		for err := range producer.Errors() {
			log.Printf("Send message failed: %s\n", err)
		}
	}()

	return p
}

func sendMsg(p sarama.AsyncProducer, topic string, value string) {
	p.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder("Hello World"),
		Value: sarama.ByteEncoder(value),
	}
}

func main() {
	tenant := flag.String("t", "ssddn", "Tenant")
	product := flag.String("p", "lock", "Product")
	flag.Parse()

	mqttBus := fmt.Sprintf(fmtOfMqttEventBus, *tenant, *product)

	sarama.Logger = logger

	producer = newProducer()

	for i := 0; i < 3; i++ {
		sendMsg(producer, mqttBus, "This is test")
		time.Sleep(2 * time.Second)
	}
	producer.Close()
	log.Println("producer over")
}
