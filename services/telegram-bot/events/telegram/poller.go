package telegram

import (
	"log"
	"time"

	"github.com/timeb30/techstreamshop/services/telegram-bot/events"
)

type BrokerListener interface {
	Consume() (*events.Event, error)
	Subscribe(topic string) error
}
type Poller struct {
	fetcher        events.Fetcher
	processor      events.Processor
	BrokerListener BrokerListener
	batchSize      int64
}

func NewPoller(fetcher events.Fetcher, processor events.Processor, batchSize int64) *Poller {
	return &Poller{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}
func (p *Poller) startListener() {
	for {
		event, err := p.BrokerListener.Consume()
		if err != nil {
			log.Println(err)
		}
		err = p.processor.Process(*event)
		if err != nil {
			log.Println(err)
		}
	}
}

func (p *Poller) Start(timeOut time.Duration) {
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

func (p *Poller) HandleEvents(e []events.Event) {
	for _, ev := range e {
		if err := p.processor.Process(ev); err != nil {
			log.Printf("failed to process event: %v", err)
			continue
		}
	}
}
