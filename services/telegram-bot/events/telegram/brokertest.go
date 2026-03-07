package telegram

import (
	"encoding/json"
	"fmt"

	"github.com/timeb30/techstreamshop/pkg/kafkaclient"
	"github.com/timeb30/techstreamshop/services/telegram-bot/events"
)

type Broker struct {
	producer *kafkaclient.Producer
	consumer *kafkaclient.Consumer
}

func newBroker(prod *kafkaclient.Producer, cons *kafkaclient.Consumer) (*Broker, error) {
	return &Broker{
		producer: prod,
		consumer: cons,
	}, nil
}
func (b *Broker) Produce(topic string, payload any) error {
	return b.producer.Produce(topic, payload)
}
func (b *Broker) Consume() (*events.Event, error) {
	msg, err := b.consumer.Consume()
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

func (b *Broker) Close() (int, error) {
	return b.producer.Close(), b.consumer.Close()
}
