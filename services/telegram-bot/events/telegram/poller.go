package telegram

import (
	"fmt"
	"log"
	"time"

	"github.com/timeb30/techstreamshop/services/telegram-bot/events"
	"github.com/timeb30/techstreamshop/services/telegram-bot/lib/e"
)

type Poller struct {
	fetcher   events.Fetcher
	processor events.Processor
	broker    events.Broker
	batchSize int64
}

func NewPoller(fetcher events.Fetcher, processor events.Processor, broker events.Broker, batchSize int64) *Poller {
	return &Poller{
		fetcher:   fetcher,
		processor: processor,
		broker:    broker,
		batchSize: batchSize,
	}
}
func (p *Poller) startConsumer() {
	const op = "Poller.startConsumer"
	for {
		event, err := p.broker.Consume()
		if err != nil {
			log.Println(e.Wrap(op, err))
			continue
		}
		err = p.processor.Process(*event)
		if err != nil {
			log.Println(e.Wrap(op, err))
		}
	}
}

func (p *Poller) StartFetcher(timeOut time.Duration) {
	fmt.Println("fetcher started")
	for {
		gotEvents, err := p.fetcher.Fetch(p.batchSize)
		if err != nil {
			log.Printf("failed to fetch events: %v", err)
			continue
		}
		if len(gotEvents) == 0 {
			time.Sleep(timeOut)
			continue
		}
		fmt.Println("got events")
		go func() {
			p.HandleEvents(gotEvents)
		}()
	}
}
func (p *Poller) Start() {
	go p.startConsumer()
	fmt.Println("Consumer started")
	p.StartFetcher(3 * time.Second)
}
func (p *Poller) HandleEvents(e []events.Event) {
	for _, ev := range e {
		if err := p.processor.Process(ev); err != nil {
			log.Printf("failed to process event: %v", err)
			continue
		}
	}
}
