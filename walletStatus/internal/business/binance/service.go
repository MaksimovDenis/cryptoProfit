package binance

import (
	"context"
	"walletStatus/internal/domain"
)

type binanceClient interface {
	GetPriceByTickers(ctx context.Context, tickers string) (map[string]float64, error)
}

type Service struct {
	binanceClient binanceClient
	stockData     map[string]domain.Stock
}

func New(
	binanceClient binanceClient,
	stockData map[string]domain.Stock,
) *Service {
	return &Service{
		binanceClient: binanceClient,
		stockData:     stockData,
	}
}
