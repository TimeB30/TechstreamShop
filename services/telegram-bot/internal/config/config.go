package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type KafkaConfig struct {
	Producer map[string]interface{} `yaml:"producer"`
	Consumer map[string]interface{} `yaml:"consumer"`
}
type Config struct {
	Env           string      `yaml:"env"`
	PostgresqlUri string      `yaml:"postgresql_uri"`
	TgBotHost     string      `yaml:"tg_bot_host"`
	KafkaConfig   KafkaConfig `yaml:"kafka"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("Error gettinng config path\n")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Error config file does not exist on: %s", configPath)
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Error parsing config file: %s", err)
	}
	return &config
}
