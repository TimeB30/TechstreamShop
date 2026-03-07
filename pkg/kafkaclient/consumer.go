package kafkaclient

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	consumer *kafka.Consumer
}

func NewConsumer(broker string) (*Consumer, error) {
	cons, err := kafka.NewConsumer(&kafka.ConfigMap{
		"boostrap.servers":  broker,
		"auto.offset.reset": "earliest",
	})
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
func (c *Consumer) Consume() (*kafka.Message, error) {
	msg, err := c.consumer.ReadMessage(time.Second)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
func (c *Consumer) Close() error {
	return c.consumer.Close()
}
