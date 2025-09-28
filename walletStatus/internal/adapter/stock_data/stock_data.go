package stockdata

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"walletStatus/internal/domain"
)

//go:embed stock-data.json
var stockData []byte

type stockDTO struct {
	Ticker       string  `json:"ticker"`
	AveragePrice float64 `json:"average_price"`
	Balance      float64 `json:"balance"`
}

type priceByTicker = map[string]domain.Stock

type StocksData struct {
	PriceByTicker priceByTicker
}

func New() (*StocksData, error) {
	var data []stockDTO
	err := json.Unmarshal(stockData, &data)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal failed to parse stock-data.json: %w", err)
	}

	priceByTicker := make(priceByTicker, len(data))

	for _, item := range data {
		priceByTicker[item.Ticker] = domain.Stock{
			Ticker:       item.Ticker,
			AveragePrice: item.AveragePrice,
			Balance:      item.Balance,
		}
	}

	return &StocksData{
		PriceByTicker: priceByTicker,
	}, nil
}
