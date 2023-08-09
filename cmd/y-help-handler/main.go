package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// upon connection to the client, this is called
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

// this is called when the connection to the client is lost, it prints "Connection lost" and the corresponding error
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

var updateMessageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	go func() {
		fmt.Println("Sending response...")
		if token := client.Publish("foobar", 1, false, "response"); token.Wait() && token.Error() != nil {
			fmt.Println("Unable to publish message...")
		}
	}()
}

var actionMessageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	go func() {
		fmt.Println("Sending response...")
		if token := client.Publish("foobar", 1, false, "response"); token.Wait() && token.Error() != nil {
			fmt.Println("Unable to publish message...")
		}
	}()
}

func main() {
	broker := "localhost"
	port := 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	// opts.SetClientID("<client_name>")
	// opts.SetUsername("<username>")
	// opts.SetPassword("<password>")

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Unable to connect to server...")
		return
	}

	if token := client.Subscribe("byu/+/update", 1, updateMessageHandler); token.Wait() && token.Error() != nil {
		fmt.Println("Unable to subscribe to topic...")
		return
	}

	if token := client.Subscribe("byu/+/+/+", 1, actionMessageHandler); token.Wait() && token.Error() != nil {
		fmt.Println("Unable to subscribe to topic...")
		return
	}

	time.Sleep(50 * time.Second)
}
