package telegram

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/timeb30/techstreamshop/services/telegram-bot/events"
)

const (
	KeysTopic = "keys"
)

type KeysMessage struct {
	ChatID int64  `json:"chat_id"`
	UserID int64  `json:"user_id"`
	Key    string `json:"key"`
}

type BrokerConsumer interface {
	Consume() (*kafka.Message, error)
	Subscribe(topic string) error
}
type BrokerEventAdapter struct {
	inner BrokerConsumer
}

func NewBrokerEventAdapter(inner BrokerConsumer) *BrokerEventAdapter {
	return &BrokerEventAdapter{
		inner: inner,
	}
}
func (b *BrokerEventAdapter) Consume() (*events.Event, error) {
	msg, err := b.inner.Consume()
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
