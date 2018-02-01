package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/yhat/scrape"

	"net/http"

	"github.com/Syfaro/telegram-bot-api"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var token string

var parseError = errors.New("Parse error")

type currency struct {
	ID    int
	Name  string
	Code  string
	Price string
}

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

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		switch update.Message.Command() {
		case "start":
			resp, err := getCurrency(10)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
			}
			var m string
			for _, v := range resp {
				m += fmt.Sprintf("|%-d|%-6s|%-3s|%-5s|\n", v.ID, v.Name, v.Code, v.Price)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, m)
			bot.Send(msg)
		case "top50":
			resp, err := getCurrency(50)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
			}
			var m string
			for _, v := range resp {
				m += fmt.Sprintf("|%-d|%-6s|%-3s|%-5s|\n", v.ID, v.Name, v.Code, v.Price)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, m)
			bot.Send(msg)
		}
	}
}

func getCurrency(n int) ([]currency, error) {
	if n == 0 {
		n = 9
	}

	var currencies []currency
	req, err := http.NewRequest("GET", "https://www.coingecko.com/en", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		return nil, errors.New("Response body is nil")
	}
	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	table, ok := scrape.Find(root, scrape.ById("gecko-table"))
	if !ok {
		return nil, parseError
	}

	tbody, ok := scrape.Find(table, scrape.ByTag(atom.Tbody))
	if !ok {
		return nil, parseError
	}

	rows := scrape.FindAll(tbody, scrape.ByTag(atom.Tr))

	for k, v := range rows {
		if k >= n {
			return currencies, nil
		}
		n, ok := scrape.Find(v, scrape.ByClass("coin-content-name"))
		if !ok {
			return nil, parseError
		}
		c, ok := scrape.Find(v, scrape.ByClass("coin-content-symbol"))
		if !ok {
			return nil, parseError
		}
		p, ok := scrape.Find(v, scrape.ByClass("currency-exchangable"))
		if !ok {
			return nil, parseError
		}

		var result currency
		result.ID = k + 1
		result.Name = scrape.Text(n)
		result.Code = scrape.Text(c)
		result.Price = scrape.Text(p)

		currencies = append(currencies, result)

	}

	return currencies, nil
}
