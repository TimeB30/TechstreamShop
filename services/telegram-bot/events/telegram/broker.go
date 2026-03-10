package telegram

import (
	"encoding/json"
	"fmt"

	"github.com/timeb30/techstreamshop/pkg/kafkaclient"
	"github.com/timeb30/techstreamshop/services/telegram-bot/events"
)

const (
	KeysTopic   = "keys"
	OrdersTopic = "orders"
)

type KeysMessage struct {
	OrderMessage *OrderMessage `json:"order_message"`
	Key          string        `json:"key"`
}

type OrderMessage struct {
	UserID     int64  `json:"user_id"`
	ChatID     int64  `json:"chat_id"`
	SoftwareID string `json:"software_id"`
	Version    int64  `json:"version"`  // index
	Duration   int64  `json:"duration"` //days

}
type KafkaBroker struct {
	producer *kafkaclient.Producer
	consumer *kafkaclient.Consumer
}

func NewBroker(prod *kafkaclient.Producer, cons *kafkaclient.Consumer) (*KafkaBroker, error) {
	return &KafkaBroker{
		producer: prod,
		consumer: cons,
	}, nil
}
func (b *KafkaBroker) Produce(topic string, payload any) error {
	return b.producer.Produce(topic, payload)
}
func (b *KafkaBroker) Consume() (*events.Event, error) {
	msg, err := b.consumer.Consume(5)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, fmt.Errorf("Message is nil")
	}
	switch *msg.TopicPartition.Topic {
	case KeysTopic:
		var keyMsg KeysMessage
		err := json.Unmarshal(msg.Value, &keyMsg)
		if err != nil {
			return nil, err
		}

		return &events.Event{
			Type: events.Key,
			Text: keyMsg.Key,
			Meta: keyMsg,
		}, nil
	default:
		return nil, fmt.Errorf("Unknown message type")
	}
}

func (b *KafkaBroker) Close() (int, error) {
	return b.producer.Close(), b.consumer.Close()
}
