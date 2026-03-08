package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/timeb30/techstreamshop/pkg/kafkaclient"
	"github.com/timeb30/techstreamshop/services/telegram-bot/client/telegram"
	processor "github.com/timeb30/techstreamshop/services/telegram-bot/events/telegram"
	"github.com/timeb30/techstreamshop/services/telegram-bot/internal/config"
	"github.com/timeb30/techstreamshop/services/telegram-bot/internal/storage/postgresql"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cnfg := config.MustLoad()
	fmt.Println(cnfg)
	tgClient := telegram.New(cnfg.TgBotHost, MustToken())
	storage, err := postgresql.New(cnfg.PostgresqlUri)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("1 stage")
	producer, err := kafkaclient.NewProducer(cnfg.KafkaConfig.Producer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("2 stage")
	consumer, err := kafkaclient.NewConsumer(cnfg.KafkaConfig.Consumer)
	if err != nil {
		log.Fatal(err)
	}
	if len(cnfg.KafkaConfig.Topics) == 0 {
		log.Fatal("topics is empty")
	}
	for _, topic := range cnfg.KafkaConfig.Topics {
		err = consumer.Subscribe(topic)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("3 stage")
	broker, err := processor.NewBroker(producer, consumer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("4 stage")
	Processor := processor.NewProcessor(&tgClient, storage, broker)
	Poller := processor.NewPoller(Processor, Processor, broker, 32)
	Poller.Start()
	log.Fatal("Error, poller stopped unexpectedly")
	//logger := setupLogger()
	//tgClient = telegram.New(token)
	// consumer.Start(fetcher, proccessor)
	// fetcher
	// processor
}

func MustToken() string {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("API token not found")
	}
	return token
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
