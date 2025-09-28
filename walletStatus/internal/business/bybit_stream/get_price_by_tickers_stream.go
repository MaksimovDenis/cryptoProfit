package bybitstream

import (
	"sync"
	"walletStatus/internal/business/utils"
	"walletStatus/internal/domain"
)

func (srv *Service) StreamPricesByBit() (<-chan map[string]float64, <-chan error) {
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

	outCh := mergeChannel(resChArr)
	mergedErrCh := mergeErrors(errChArr, errCh)

	return outCh, mergedErrCh
}

func mergeChannel(channels []<-chan map[string]float64) chan map[string]float64 {
	var wg sync.WaitGroup
	outputCh := make(chan map[string]float64)

	wg.Add(len(channels))
	for _, ch := range channels {
		go func(c <-chan map[string]float64) {
			defer wg.Done()
			for value := range c {
				outputCh <- value
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()

	return outputCh
}

func mergeErrors(channels []<-chan error, errCh chan error) chan error {
	var wg sync.WaitGroup

	wg.Add(len(channels))
	for _, ch := range channels {
		go func(c <-chan error) {
			defer wg.Done()
			for err := range c {
				errCh <- err
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	return errCh
}
