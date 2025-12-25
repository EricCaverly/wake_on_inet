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

	tok := client.Subscribe(cfg.CommandTopic, byte(cfg.CommandQOS), cmd_handler)
	if tok.Wait() && tok.Error() != nil {
		return client, tok.Error()
	}

	return client, nil
}
