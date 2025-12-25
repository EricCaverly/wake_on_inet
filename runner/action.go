package main

import (
	"encoding/json"
	"log"
	"slices"

	"github.com/EricCaverly/wake_on_inet/common"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	valid_subnets []string
)

var wake_cmd_handler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	wake_cmd := common.WakeCommand{}

	err := json.Unmarshal(msg.Payload(), &wake_cmd)
	if err != nil {
		log.Printf("Malformed wake command: %s\n", err.Error())
		return
	}

	if !slices.Contains(valid_subnets, wake_cmd.Subnet) {
		log.Printf("Subnet '%s' not valid for this runner\n", wake_cmd.Subnet)
		return
	}

	err = wake_pc(wake_cmd.MacAddress, wake_cmd.Subnet)
	if err != nil {
		log.Printf("Failed to send WakeOnLan packet: %s\n", err.Error())
		return
	}

	log.Printf("Processed wake command for %s on %s\n", wake_cmd.MacAddress, wake_cmd.Subnet)
}

var ping_cmd_handler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	ping_cmd := common.PingCommand{}

	err := json.Unmarshal(msg.Payload(), &ping_cmd)
	if err != nil {
		log.Printf("Malformed ping command: %s\n", err.Error())
	}

	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}
