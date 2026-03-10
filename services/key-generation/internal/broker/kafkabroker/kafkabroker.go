package kafkabroker

import (
	"encoding/json"
	"fmt"

	"github.com/timeb30/techstreamshop/pkg/kafkaclient"
	"github.com/timeb30/techstreamshop/services/key-generation/internal/broker"
)

const (
	KeysTopic   = "keys"
	OrdersTopic = "orders"
)

type KafkaBroker struct {
	producer *kafkaclient.Producer
	consumer *kafkaclient.Consumer
	//TODO add redis to make consumer idempotent
}

func NewBroker(prod *kafkaclient.Producer, cons *kafkaclient.Consumer) (*KafkaBroker, error) {
	return &KafkaBroker{
		producer: prod,
		consumer: cons,
	}, nil
}
func (b *KafkaBroker) Produce(message *broker.KeyMessage) error {
	return b.producer.Produce(KeysTopic, message)
}
func (b *KafkaBroker) Consume(timeOutDuration int) (*broker.OrderMessage, error) {
	msg, err := b.consumer.Consume(timeOutDuration)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, fmt.Errorf("Message is nil")
	}
	switch *msg.TopicPartition.Topic {
	case OrdersTopic:
		var orderMsg broker.OrderMessage
		err := json.Unmarshal(msg.Value, &orderMsg)
		if err != nil {
			return nil, err
		}
		return &orderMsg, nil
	default:
		return nil, fmt.Errorf("Unknown topic")
	}
}

func (b *KafkaBroker) Close() (int, error) {
	return b.producer.Close(), b.consumer.Close()
}
