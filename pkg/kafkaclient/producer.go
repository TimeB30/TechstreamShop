package kafkaclient

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(config map[string]interface{}) (*Producer, error) {
	kafkaConfig := getKafkaConfig(config)
	producer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		return nil, err
	}
	return &Producer{producer: producer}, nil
}

func (p *Producer) Produce(topic string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: data,
	}, deliveryChan)
	if err != nil {
		return err
	}
	e := <-deliveryChan
	m, ok := e.(*kafka.Message)
	if !ok {
		return fmt.Errorf("unexpected event delivered from kafka producer")
	}
	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	} else {
		return nil
	}
}
func (p *Producer) Close() int {
	count := p.producer.Flush(15 * 1000)
	p.producer.Close()
	return count
}
