package telegram

import (
	"fmt"
	"strconv"

	"github.com/timeb30/techstreamshop/services/telegram-bot/client/telegram"
	"github.com/timeb30/techstreamshop/services/telegram-bot/lib/e"
)

const (
	HelpCmd  = "/help"
	StartCmd = "/start"
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
			var replyMarkup telegram.InlineKeyboardMarkup
			leftIndex := 0
			rightIndex := inlineKeyboardMarkupLen - 1
			replyMarkup, err := p.getKeyboard(leftIndex, rightIndex, SoftwareVersions, replyMarkup, fmt.Sprintf("%s/%s", SetVersionQuery, text))
			replyMarkup = append(replyMarkup, p.getButton("⬇️", fmt.Sprintf("%s/%s", VersionDownQuery, text), strconv.Itoa(leftIndex), strconv.Itoa(rightIndex)))
			if err != nil {
				return e.Wrap("can't attach markup", err)
			}
			err = p.tg.SendMessage(chatID, text, &telegram.ReplyMarkup{InlineKeyboardMarkup: replyMarkup})
			if err != nil {
				return err
			}

		} else {
			err := p.tg.SendMessage(chatID, unknownCommand, nil)
			if err != nil {
				return err
			}
			return fmt.Errorf("Software ID should be 32 symbols")
		}
	}
	return nil
}
