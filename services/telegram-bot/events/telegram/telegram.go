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
	case events.CallBackQuery:
		return p.processQuery(e)
	default:
		return ErrUnknownEvent
	}
}
func (p *Processor) processQuery(e events.Event) error {
	res, ok := e.Meta.(events.CallBackQueryMeta)
	if !ok || res.Data == "" {
		return ErrUnknownMetaType
	}
	return p.doQuery(res.CallbackQueryID, res.ChatID, res.UserID, res.MessageID, res.Data)
}
func (p *Processor) processKey(e events.Event) error {
	res, ok := e.Meta.(KeysMessage)
	if !ok {
		return ErrUnknownMetaType
	}
	msg := fmt.Sprintf("%s\n%s\nВерсия: %s\nСрок действия(дни): %d", res.OrderMessage.SoftwareID, e.Text, SoftwareVersions[res.OrderMessage.Version], res.OrderMessage.Duration)
	err := p.tg.SendMessage(res.OrderMessage.ChatID, msg, nil)
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
	switch updType {
	case events.Message:
		res.Meta = events.MessageMeta{
			ChatID:   upd.Message.Chat.ChatID,
			UserName: upd.Message.From.UserName,
			UserID:   upd.Message.From.UserID,
		}
	case events.CallBackQuery:
		res.Meta = events.CallBackQueryMeta{
			UserID:          upd.CallBackQuery.From.UserID,
			ChatID:          upd.CallBackQuery.Message.Chat.ChatID,
			CallbackQueryID: upd.CallBackQuery.ID,
			MessageID:       upd.CallBackQuery.Message.ID,
			Data:            upd.CallBackQuery.Data,
		}
	default:
		res.Meta = nil
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
	if upd.CallBackQuery != nil {
		return events.CallBackQuery
	}
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}
