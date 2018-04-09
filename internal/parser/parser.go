package parser

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var errEmptyCoinName = errors.New("coin name is empty")

type Currency struct {
	Name       string `json:"name"`
	Code       string `json:"symbol"`
	MarketData struct {
		Price struct {
			USD Number `json:"usd"`
		} `json:"current_price"`
	} `json:"market_data"`
}

type Parser struct {
	client *http.Client
}

func NewClient() Parser {
	return Parser{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (p Parser) GetCurrencyByName(name string) ([]Currency, error) {
	if name == "" {
		return nil, errEmptyCoinName
	}

	resp, err := p.client.Get("https://api.coingecko.com/api/v3/coins/" + name)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	var out Currency
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return []Currency{out}, nil
}

func (p Parser) GetCurrencyList(n string) ([]Currency, error) {
	resp, err := p.client.Get("https://api.coingecko.com/api/v3/coins?per_page=" + n)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	var out []Currency
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}
