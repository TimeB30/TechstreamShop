package telegram

import (
	"errors"
	"log/slog"

	"github.com/timeb30/techstreamshop/services/telegram-bot/client/telegram"
	"github.com/timeb30/techstreamshop/services/telegram-bot/events"
	"github.com/timeb30/techstreamshop/services/telegram-bot/internal/storage"
	"github.com/timeb30/techstreamshop/services/telegram-bot/lib/e"
)

var ErrUnknownEvent = errors.New("unknown event type")

type Processor struct {
	tg      *telegram.Client
	offset  int64
	storage *storage.Storage
	logger  *slog.Logger
}
type Meta struct {
	ChatID   int64
	UserName string
}

func NewProcessor(client *telegram.Client, storage *storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int64) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.logger, p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get updates", err)
	}
	if len(updates) == 0 {
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
	default:
		return ErrUnknownEvent
	}
}
func (p *Processor) processMessage(e events.Event) error {
	//TODO sending to key generate service
}
func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ChatID,
			UserName: upd.Message.From.UserName,
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
