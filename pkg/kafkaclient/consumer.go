package kafkaclient

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	consumer *kafka.Consumer
}

func NewConsumer(config map[string]interface{}) (*Consumer, error) {
	kafkaConfig := getKafkaConfig(config)
	cons, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		consumer: cons,
	}, nil
}
func (c *Consumer) Subscribe(topic string) error {
	return c.consumer.Subscribe(topic, nil)
}
func (c *Consumer) Consume(timeOutSec int) (*kafka.Message, error) {
	msg, err := c.consumer.ReadMessage(time.Duration(timeOutSec) * time.Second)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
func (c *Consumer) Close() error {
	return c.consumer.Close()
}
