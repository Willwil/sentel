package main

import (
	"flag"
	"fmt"
	//import the Paho Go MQTT library
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s, MSG: %s", msg.Topic(), msg.Payload())
}

func main() {
	clientID := flag.String("c", "mqtt-client", "Client ID")
	topic := flag.String("t", "device1/product1/data", "Topic")
	port := flag.Int("p", 1883, "Port")
	bePublish := flag.Bool("pub", false, "Publish")
	beSubscribe := flag.Bool("sub", true, "Subscribe")
	flag.Parse()

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker(fmt.Sprintf("tcp://localhost:%d", *port))
	opts.SetClientID(fmt.Sprintf("%s|securemode=3,signmethod=hmacsha1,timestamp=12345", *clientID))
	opts.SetUsername("device1&product2")
	opts.SetPassword("123456")
	opts.SetDefaultPublishHandler(f)
	opts.SetCleanSession(true)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Mqtt Connect success")

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if *beSubscribe {
		if token := c.Subscribe(*topic, 0, nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
		fmt.Println("Mqtt SUBSCRIBE success")
		time.Sleep(3 * time.Second)
	}

	if *bePublish {
		for i := 0; i < 2; i++ {
			text := fmt.Sprintf("this is msg #%d!", i)
			token := c.Publish(*topic, 0, false, text)
			token.Wait()
			if token.Error() != nil {
				fmt.Println("Failed to PUBLISH message to topic")
			} else {
				fmt.Println("Mqtt PUBLISH success")
			}

			time.Sleep(1 * time.Second)
		}
	}

	time.Sleep(100 * time.Second)

	if *beSubscribe {
		//unsubscribe from /go-mqtt/sample
		if token := c.Unsubscribe(*topic); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
	}

	c.Disconnect(250)
}
