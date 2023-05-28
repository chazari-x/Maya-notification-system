package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

var C Config

type Config struct {
	RabbitMQ string `yaml:"rabbitMQ"`

	Bot struct {
		Token string `yaml:"token"`
		URL   string `yaml:"url"`
	} `yaml:"bot"`
}

const configFile = "worker/internal/app/config/dev.yaml"

func GetConfig() (*Config, error) {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return nil, errors.New("read config fIle err: " + err.Error())
	}

	if err = yaml.Unmarshal(yamlFile, &C); err != nil {
		return nil, errors.New("unmarshal config fIle err: " + err.Error())
	}

	if C.RabbitMQ == "" || C.Bot.URL == "" || C.Bot.Token == "" {
		return nil, errors.New("error config")
	}

	return &C, nil
}
