package main

import (
	"fmt"
	"github.com/timeb30/techstreamshop/pkg/kafkaclient"
	"github.com/timeb30/techstreamshop/services/key-generation/internal/broker/kafkabroker"
	"github.com/timeb30/techstreamshop/services/key-generation/internal/config"
	"github.com/timeb30/techstreamshop/services/key-generation/internal/keygen"
	"github.com/timeb30/techstreamshop/services/key-generation/internal/processor"
	"github.com/timeb30/techstreamshop/services/key-generation/internal/storage/postgresql"
	"log"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cnfg := config.MustLoad()
	producer, err := kafkaclient.NewProducer(cnfg.KafkaConfig.Producer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("2 stage")
	consumer, err := kafkaclient.NewConsumer(cnfg.KafkaConfig.Consumer)
	if err != nil {
		log.Fatal(err)
	}
	err = consumer.Subscribe(cnfg.KafkaConfig.Topics[0])
	fmt.Println("Subscribed to", cnfg.KafkaConfig.Topics[0])
	if err != nil {
		log.Fatal(err)
	}
	brkr, err := kafkabroker.NewBroker(producer, consumer)
	if err != nil {
		log.Fatal(err)
	}
	storage, err := postgresql.New(cnfg.PostgresqlUri)
	if err != nil {
		log.Fatal(err)
	}
	kg := keygen.NewKeyGen("http://host.docker.internal:8000", nil)
	procsr := processor.NewProcessor(brkr, *kg, *storage)
	procsr.Start()
	log.Fatal("Processor stop unexpectedly")

}
