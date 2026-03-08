package telegram

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/timeb30/techstreamshop/services/telegram-bot/client/telegram"
	"github.com/timeb30/techstreamshop/services/telegram-bot/events"
	"github.com/timeb30/techstreamshop/services/telegram-bot/internal/storage"
	"github.com/timeb30/techstreamshop/services/telegram-bot/lib/e"
)

var (
	ErrUnknownEvent    = errors.New("unknown event type")
	ErrUnknownMetaType = errors.New("unknown meta type")
)

type Processor struct {
	tg      *telegram.Client
	offset  int64
	storage storage.Storage
	logger  *slog.Logger
	broker  events.Broker
}

func NewProcessor(client *telegram.Client, storage storage.Storage, broker events.Broker) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
		broker:  broker,
	}
}

func (p *Processor) Fetch(limit int64) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get updates", err)
	}
	if len(updates) == 0 {
		fmt.Println("No updates")
		return nil, nil
	}
	res := make([]events.Event, 0, len(updates))
	for _, u := range updates {
		res = append(res, event(u))
	}
	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func (p *Processor) Process(e events.Event) error {
	switch e.Type {
	case events.Message:
		return p.processMessage(e)
	case events.Key:
		return p.processKey(e)
	default:
		return ErrUnknownEvent
	}
}
func (p *Processor) processKey(e events.Event) error {
	res, ok := e.Meta.(events.KeyMeta)
	if !ok {
		return ErrUnknownMetaType
	}
	err := p.tg.SendMessage(res.ChatID, e.Text, nil)
	if err != nil {
		return err
	}
	return nil
}
func (p *Processor) processMessage(e events.Event) error {
	res, ok := e.Meta.(events.MessageMeta)
	if !ok {
		return ErrUnknownMetaType
	}
	err := p.doCmd(e.Text, res.ChatID, res.UserID)
	if err != nil {
		return err
	}
	return nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message {
		res.Meta = events.MessageMeta{
			ChatID:   upd.Message.Chat.ChatID,
			UserName: upd.Message.From.UserName,
			UserID:   upd.Message.From.UserID,
		}
	}
	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}
