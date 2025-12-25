package main

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var cmd_handler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}
