package binance

import (
	"context"
	"fmt"
	"math"
	"walletStatus/internal/domain"
)

func (srv *Service) GetProfitByTickers(ctx context.Context) (map[string]domain.Profit, error) {
	currentPrice, err := srv.GetPriceByTickers(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetPriceByTickers: %w", err)
	}

	result, err := srv.calculateProfit(currentPrice, srv.stockData)
	if err != nil {
		return nil, fmt.Errorf("srv.calculateProfit: %w", err)
	}

	return result, nil
}

func (srv *Service) calculateProfit(currentPrice map[string]float64, stockData map[string]domain.Stock) (map[string]domain.Profit, error) {

	result := make(map[string]domain.Profit, len(stockData))

	for symbol, currentPrice := range currentPrice {
		stock, ok := stockData[symbol]
		if !ok {
			return nil, domain.ErrPricesCount
		}

		profitPercent := (currentPrice - stock.AveragePrice) / (stock.AveragePrice) * 100
		profit := profitPercent * (stock.Balance * currentPrice) / 100
		balance := currentPrice * stock.Balance

		result[symbol] = domain.Profit{
			Profit:        roundFloat(profit, 4),
			ProfitPercent: roundFloat(profitPercent, 4),
			Balance:       roundFloat(balance, 4),
		}
	}

	return result, nil
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
