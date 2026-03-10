package telegram

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/timeb30/techstreamshop/services/telegram-bot/client/telegram"
)

const (
	VersionDownQuery = "versionDown"
	VersionUpQuery   = "versionUp"
	SetDurationQuery = "Duration"
	SetVersionQuery  = "Version"
)

func (p *Processor) doQuery(callBackQueryID string, chatID int64, userID int64, messageID int64, data string) error {
	params := strings.Split(data, "/")
	query := params[0]
	softwareID := params[1]
	switch query {
	case VersionDownQuery, VersionUpQuery:
		reply := make(telegram.InlineKeyboardMarkup, 0, inlineKeyboardMarkupLen+2)
		leftIndex, _ := strconv.Atoi(params[2])
		rightIndex, _ := strconv.Atoi(params[3])
		lastIndex := len(SoftwareVersions) - 1
		sign := 1
		if query == VersionUpQuery {
			sign = -1
		}
		leftIndex += sign * inlineKeyboardMarkupLen
		rightIndex += sign * inlineKeyboardMarkupLen
		if leftIndex <= 0 {
			leftIndex = 0
		} else {
			reply = append(reply, p.getButton("⬆️", fmt.Sprintf("%s/%s", VersionUpQuery, softwareID), strconv.Itoa(leftIndex), strconv.Itoa(rightIndex)))
		}
		if rightIndex > lastIndex {
			rightIndex = lastIndex
		}
		reply, err := p.getKeyboard(leftIndex, rightIndex, SoftwareVersions, reply, fmt.Sprintf("%s/%s", SetVersionQuery, softwareID))
		if err != nil {
			return err
		}
		if rightIndex < lastIndex {
			reply = append(reply, p.getButton("⬇️", fmt.Sprintf("%s/%s", VersionDownQuery, softwareID), strconv.Itoa(leftIndex), strconv.Itoa(rightIndex)))
		}
		p.tg.AnswerCallBackQuery(callBackQueryID, "Select version")
		res, err := p.tg.EditMessageReplyMarkup(messageID, chatID, &telegram.ReplyMarkup{InlineKeyboardMarkup: reply})
		if err != nil {
			return err
		}
		if !res {
			return fmt.Errorf("MarkUp Edit Error")
		}

	case SetVersionQuery:
		version := params[2]
		p.tg.AnswerCallBackQuery(callBackQueryID, "Select duration")
		var reply telegram.InlineKeyboardMarkup
		reply, err := p.getKeyboard(0, len(Durations)-1, Durations, reply, fmt.Sprintf("%s/%s", fmt.Sprintf("%s/%s", SetDurationQuery, softwareID), version))
		if err != nil {
			return err
		}
		res, err := p.tg.EditMessageReplyMarkup(messageID, chatID, &telegram.ReplyMarkup{InlineKeyboardMarkup: reply})
		if err != nil {
			return err
		}
		if !res {
			return fmt.Errorf("MarkUp Edit Error")
		}
	case SetDurationQuery:
		fmt.Println(data)
		softwareVersionIndex, _ := strconv.ParseInt(params[2], 10, 64)
		durationIndex, _ := strconv.Atoi(params[3])
		duration, _ := strconv.ParseInt(Durations[durationIndex], 10, 64)
		err := p.broker.Produce(OrdersTopic, OrderMessage{
			UserID:     userID,
			ChatID:     chatID,
			SoftwareID: softwareID,
			Version:    softwareVersionIndex,
			Duration:   duration,
		})
		if err != nil {
			return err
		}
		_ = p.tg.DeleteMessage(chatID, messageID)
	default:
		fmt.Println("Unknown query")
	}

	return nil
}
