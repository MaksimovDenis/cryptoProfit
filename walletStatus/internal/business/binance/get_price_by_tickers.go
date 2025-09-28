package binance

import (
	"context"
	"fmt"
	"walletStatus/internal/business/utils"
	"walletStatus/internal/domain"
)

func (srv *Service) GetPriceByTickers(ctx context.Context) (map[string]float64, error) {
	tickers := utils.ConvTickersToStrBinance(srv.stockData)
	if tickers == "" {
		return nil, domain.ErrTickersNotFound
	}

	prices, err := srv.binanceClient.GetPriceByTickers(ctx, tickers)
	if err != nil {
		return nil, fmt.Errorf("binanceClient.GetPriceByTicker: %w", err)
	}

	return prices, nil
}
