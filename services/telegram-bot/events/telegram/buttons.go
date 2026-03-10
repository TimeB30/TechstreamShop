package telegram

import (
	"fmt"
	"strconv"

	"github.com/timeb30/techstreamshop/services/telegram-bot/client/telegram"
)

const (
	inlineKeyboardMarkupLen = 8
)

var SoftwareVersions = []string{
	"18.00.008",
	"17.30.011",
	"17.20.013",
	"17.10.012",
	"17.00.020",
	"16.30.013",
	"16.20.026",
	"16.20.023",
	"16.10.016",
	"16.00.017",
	"15.30.027",
	"15.30.026",
	"15.20.015",
	"15.10.029",
	"15.00.028",
	"15.00.026",
	"14.30.023",
	"14.20.019",
	"14.10.033",
	"14.10.028",
	"14.00.019",
	"14.00.018",
	"13.30.018",
	"13.20.018",
	"13.20.017",
	"13.10.019",
	"13.00.022",
	"12.30.017",
	"12.20.024",
	"12.10.019",
	"12.00.127",
	"12.00.125",
	"12.00.124",
	"11.30.137",
	"11.30.124",
	"11.30.037",
	"11.30.024",
	"11.20.019",
	"11.10.034",
	"11.00.019",
	"11.00.017",
}

var Durations = []string{
	"5000",
	"365",
	"60",
	"7",
	"1",
}

func (p *Processor) getButton(text string, command string, options ...string) []telegram.InlineKeyboardButton {
	optionsStr := ""
	for _, option := range options {
		optionsStr += "/" + option
	}
	return []telegram.InlineKeyboardButton{
		{
			Text:         text,
			CallBackData: fmt.Sprintf("%s%s", command, optionsStr),
		},
	}
}

func (p *Processor) getKeyboard(startIndex int, endIndex int, btnTexts []string, markUp telegram.InlineKeyboardMarkup, queryType string) (telegram.InlineKeyboardMarkup, error) {
	if startIndex < 0 || endIndex > len(btnTexts)-1 || startIndex > endIndex {
		return nil, fmt.Errorf("Wrong index")
	}
	for i := startIndex; i <= endIndex; i++ {
		markUp = append(markUp, p.getButton(btnTexts[i], queryType, strconv.Itoa(i)))
	}
	return markUp, nil
}
