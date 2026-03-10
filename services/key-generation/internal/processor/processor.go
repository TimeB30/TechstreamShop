package processor

import (
	"fmt"
	"log"
	"time"

	"github.com/timeb30/techstreamshop/services/key-generation/internal/broker"
	"github.com/timeb30/techstreamshop/services/key-generation/internal/keygen"
	"github.com/timeb30/techstreamshop/services/key-generation/internal/storage/postgresql"
)

type Processor struct {
	broker  broker.Broker
	keygen  keygen.KeyGenerator
	storage postgresql.Storage
}

func NewProcessor(broker broker.Broker, generator keygen.KeyGenerator, storage postgresql.Storage) *Processor {
	return &Processor{
		broker:  broker,
		keygen:  generator,
		storage: storage,
	}
}

func (p *Processor) ProcessOrder(orderMsg *broker.OrderMessage) error {
	key, err := p.keygen.GenerateKey(orderMsg.SoftwareID, orderMsg.Duration, orderMsg.Version)
	if err != nil {
		return err
	}
	fmt.Println("Key generated:", key)
	_, err = p.storage.SaveKey(orderMsg.UserID, key, time.Now(), time.Now().Add(time.Duration(orderMsg.Duration)*24*time.Hour))
	if err != nil {
		return err
	}
	res := broker.KeyMessage{
		OrderMessage: orderMsg,
		Key:          key,
	}
	err = p.broker.Produce(&res)
	if err != nil {
		return err
	}
	return nil
}

func (p *Processor) Start() {
	for {
		orderMsg, err := p.broker.Consume(5)
		log.Println("Consuming")
		if err != nil {
			log.Println("Error consuming order:", err)
			continue
		}
		if orderMsg == nil {
			log.Println("Error consuming order: nil orderMsg")
			continue
		}
		go func() {
			err := p.ProcessOrder(orderMsg)
			if err != nil {
				log.Println("Error processing order:", err)
			}
		}()

	}
}
