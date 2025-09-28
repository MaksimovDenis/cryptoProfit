package bybitstream

import (
	"walletStatus/internal/domain"
)

type ByBitStream interface {
	SubscribeTickers(tickers string) (<-chan map[string]float64, chan error)
}

type Service struct {
	stockData   map[string]domain.Stock
	byBitStream ByBitStream
}

func New(byBitStream ByBitStream, stockData map[string]domain.Stock) *Service {
	return &Service{
		byBitStream: byBitStream,
		stockData:   stockData,
	}
}
