package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type KafkaConfig struct {
	Topics   []string               `yaml:"topics" env-required`
	Producer map[string]interface{} `yaml:"producer" env-required`
	Consumer map[string]interface{} `yaml:"consumer" env-required`
}
type Config struct {
	Env           string      `yaml:"env" env-required`
	PostgresqlUri string      `yaml:"postgresql_uri" env-required`
	TgBotHost     string      `yaml:"tg_bot_host" env-required`
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

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("Error parsing config file: %s", err)
	}
	config.PostgresqlUri = os.ExpandEnv(config.PostgresqlUri)
	return &config
}
