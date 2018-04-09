package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bynov/CryptoCurrencyBot/internal/parser"

	"github.com/Syfaro/telegram-bot-api"
)

var token string

const PREFIX = "getfucking"

func init() {
	if v, ok := os.LookupEnv("CRYPTOCURRENCY_TOKEN"); ok {
		token = v
	} else {
		panic("CRYPTOCURRENCY_TOKEN is not provided")
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	requestor := parser.NewClient()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		cmd := strings.ToLower(update.Message.Command())
		switch cmd {
		case "start":
			resp, err := requestor.GetCurrencyList("10")
			if err != nil {
				if _, err := reportError(bot, update.Message.Chat.ID, err); err != nil {
					log.Println(err)
				}

				continue
			}

			if _, err := sendResponse(bot, update.Message.Chat.ID, resp); err != nil {
				log.Println(err)
			}
		case "top50":
			resp, err := requestor.GetCurrencyList("50")
			if err != nil {
				if _, err := reportError(bot, update.Message.Chat.ID, err); err != nil {
					log.Println(err)
				}

				continue
			}

			if _, err := sendResponse(bot, update.Message.Chat.ID, resp); err != nil {
				log.Println(err)
			}
		default:
			if !strings.Contains(cmd, PREFIX) {
				return
			}
			t := strings.TrimPrefix(cmd, PREFIX)
			resp, err := requestor.GetCurrencyByName(t)
			if err != nil {
				if _, err := reportError(bot, update.Message.Chat.ID, err); err != nil {
					log.Println(err)
				}

				continue
			}

			if _, err := sendResponse(bot, update.Message.Chat.ID, resp); err != nil {
				log.Println(err)
			}
		}
	}
}

func reportError(bot *tgbotapi.BotAPI, chatID int64, err error) (tgbotapi.Message, error) {
	return bot.Send(tgbotapi.NewMessage(chatID, err.Error()))
}

func sendResponse(bot *tgbotapi.BotAPI, chatID int64, currencies []parser.Currency) (tgbotapi.Message, error) {
	var m string
	var id = 1
	for _, v := range currencies {
		m += fmt.Sprintf("|%-d|%-6s|%-3s|%-5.2f$|\n", id, v.Name, v.Code, v.MarketData.Price.USD)
		id++
	}

	msg := tgbotapi.NewMessage(chatID, m)
	return bot.Send(msg)
}
