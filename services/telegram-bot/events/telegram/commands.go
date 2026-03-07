package telegram

import (
	"fmt"
	"github.com/google/uuid"
)

const (
	HelpCmd  = "/help"
	StartCmd = "/start"
)
const (
	OrdersTopic = "orders"
)

func (p *Processor) doCmd(text string, chatID int64, userID int64) error {
	switch text {
	case HelpCmd, StartCmd:
		err := p.tg.SendMessage(chatID, helpMessage, nil)
		if err != nil {
			return err
		}
	default:
		if len(text) == 32 {
			err := p.brokerProducer.Produce(OrdersTopic, struct {
				Key        string `json:"key"`
				UserID     int64  `json:"user_id"`
				SoftwareID string `json:"software_id"`
			}{
				Key:    uuid.New().String(),
				UserID: userID,
			})
			if err != nil {
				return err
			}
		}
		err := p.tg.SendMessage(chatID, "Software ID should be 32 symbols", nil)
		if err != nil {
			return err
		}
		return fmt.Errorf("Software ID should be 32 symbols")
	}
	return nil
}
