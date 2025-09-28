package bybitstream

import (
	"math"
	"sync"
	"walletStatus/internal/business/utils"
	"walletStatus/internal/domain"
)

func (srv *Service) StreamProfitByBit() (chan map[string]domain.Profit, <-chan error) {
	errCh := make(chan error, 10)

	tickers := utils.ConvTickersToStrByBit(srv.stockData)
	if len(tickers) == 0 {
		go func() { errCh <- domain.ErrTickersNotFound }()
		return nil, errCh
	}

	resChArr := make([]<-chan map[string]float64, len(tickers))
	errChArr := make([]<-chan error, len(tickers))

	for idx, ticker := range tickers {
		resChArr[idx], errChArr[idx] = srv.byBitStream.SubscribeTickers(ticker)
	}

	outCh := srv.processAndMergeChannels(resChArr)
	mergedErrCh := mergeErrors(errChArr, errCh)

	return outCh, mergedErrCh
}

func (srv *Service) processAndMergeChannels(channels []<-chan map[string]float64) chan map[string]domain.Profit {
	var wg sync.WaitGroup
	outputCh := make(chan map[string]domain.Profit)

	wg.Add(len(channels))
	for _, ch := range channels {
		go func(c <-chan map[string]float64) {
			defer wg.Done()
			for value := range c {
				result, _ := srv.calculateProfit(value, srv.stockData)

				outputCh <- result
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()

	return outputCh
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
