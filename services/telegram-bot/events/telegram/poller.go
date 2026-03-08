package telegram

import (
	"log"
	"time"

	"github.com/timeb30/techstreamshop/services/telegram-bot/events"
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
	for {
		event, err := p.broker.Consume()
		if err != nil {
			log.Println(err)
		}
		err = p.processor.Process(*event)
		if err != nil {
			log.Println(err)
		}
	}
}

func (p *Poller) StartFetcher(timeOut time.Duration) {
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
		go func() {
			p.HandleEvents(gotEvents)
		}()
	}
}
func (p *Poller) Start() {
	go p.startConsumer()
	p.StartFetcher(10 * time.Second)
}
func (p *Poller) HandleEvents(e []events.Event) {
	for _, ev := range e {
		if err := p.processor.Process(ev); err != nil {
			log.Printf("failed to process event: %v", err)
			continue
		}
	}
}
