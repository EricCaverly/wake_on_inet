package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var runner_exit_chan chan error

func main() {
	cfg, err := load_cfg("./config.yml")
	if err != nil {
		panic(err)
	}

	valid_subnets = cfg.Subnets

	c_sig := make(chan os.Signal, 1)
	signal.Notify(c_sig, os.Interrupt, syscall.SIGTERM)
	runner_exit_chan = make(chan error)

	for {
		client, err := work(cfg)
		if err != nil {
			log.Printf("Error connecting: %v\n", err)
		}

		log.Println("Waiting for interrupt...")

		select {
		case <-c_sig:
			client.Disconnect(200)
			return
		case err := <-runner_exit_chan:
			client.Disconnect(200)
			log.Printf("Error: %s\nwaiting 2 seconds and reconnecting\n", err.Error())
			time.Sleep(2 * time.Second)
		}
	}
}

var conn_handler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected!")
}

var disco_handler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
	runner_exit_chan <- err
}

func sub(client mqtt.Client, topic string, qos byte, handler mqtt.MessageHandler) error {
	wake_tok := client.Subscribe(topic, qos, wake_cmd_handler)
	if wake_tok.Wait() && wake_tok.Error() != nil {
		wake_tok.Error()
	}
	log.Printf("Subscribed to %s\n", topic)
	return nil
}

func work(cfg Config) (mqtt.Client, error) {

	// https://dev.to/emqx/how-to-use-mqtt-in-golang-2oek

	log.Println("Starting worker thread:")
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Broker)
	opts.SetClientID(cfg.ClientID)
	opts.SetUsername(cfg.Username)
	opts.SetPassword(cfg.password)

	opts.OnConnect = conn_handler
	opts.OnConnectionLost = disco_handler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return client, token.Error()
	}

	if err := sub(client, cfg.WakeCommandTopic, byte(cfg.CommandQOS), wake_cmd_handler); err != nil {
		return client, err
	}
	if err := sub(client, cfg.PingCommandTopic, byte(cfg.CommandQOS), ping_cmd_handler); err != nil {
		return client, err
	}

	return client, nil
}
