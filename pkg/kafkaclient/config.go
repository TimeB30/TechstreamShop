package kafkaclient

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

func getKafkaConfig(config map[string]interface{}) *kafka.ConfigMap {
	kafkaConfig := kafka.ConfigMap{}
	for key, value := range config {
		kafkaConfig[key] = value
	}
	return &kafkaConfig
}
