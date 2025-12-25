package main

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Broker        string `yaml:"broker"`
	ClientID      string `yaml:"client_id"`
	Username      string `yaml:"username"`
	password      string
	PasswordFile  string `yaml:"password_file"`
	CommandTopic  string `yaml:"command_topic"`
	CommandQOS    int    `yaml:"qos"`
	ResponseTopic string `yaml:"response_topic"`
}

func load_cfg(path string) (Config, error) {
	var cfg Config

	contents, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(contents, &cfg)
	if err != nil {
		return cfg, err
	}

	pw, err := os.ReadFile(cfg.PasswordFile)
	if err != nil {
		return cfg, err
	}

	cfg.password = strings.TrimSpace(string(pw))

	return cfg, err
}
