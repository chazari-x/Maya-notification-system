package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

var C Config

type Config struct {
	RunAddress string `yaml:"runAddress"`
	RabbitMQ   string `yaml:"rabbitMQ"`
}

const configFile = "internal/app/config/dev.yaml"

func GetConfig() (Config, error) {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}, errors.New("read config fIle err: " + err.Error())
	}

	if err = yaml.Unmarshal(yamlFile, &C); err != nil {
		return Config{}, errors.New("unmarshal config fIle err: " + err.Error())
	}

	if C.RunAddress == "" || C.RabbitMQ == "" {
		return Config{}, errors.New("error config")
	}

	return C, nil
}
