package app

import (
	"context"
	"net/http"
	bybitws "walletStatus/internal/adapter/bybit_ws"
	"walletStatus/internal/adapter/client/binance"
	stockdata "walletStatus/internal/adapter/stock_data"
	api "walletStatus/internal/api/http/handler"
	wsserver "walletStatus/internal/api/ws"
	binanceService "walletStatus/internal/business/binance"
	bybitService "walletStatus/internal/business/bybit_stream"
	"walletStatus/internal/infra/config"
	"walletStatus/internal/infra/logger"
)

type serviceProvider struct {
	config config.Config

	stockData *stockdata.StocksData

	binanceClient     *binance.Client
	httpBinanceClient *http.Client
	wsByBitClient     *bybitws.Client

	binanceService *binanceService.Service
	byBitService   *bybitService.Service

	appServer *api.Server
	wsServer  *wsserver.Server
}

func newServiceProvider() *serviceProvider {
	srv := &serviceProvider{}

	return srv
}

func (srv *serviceProvider) AppStockData(ctx context.Context) *stockdata.StocksData {
	if srv.stockData == nil {
		var err error

		srv.stockData, err = stockdata.New()
		if err != nil {
			logger.Fatalf(ctx, "AppStockData: %w", err)
		}
	}

	return srv.stockData
}

func (srv *serviceProvider) AppBinanceClient(_ context.Context) *binance.Client {
	if srv.binanceClient == nil {
		srv.binanceClient = binance.New(
			*srv.httpBinanceClient,
			srv.config.Binance.Address,
		)
	}

	return srv.binanceClient
}

func (srv *serviceProvider) AppBinanceService(ctx context.Context) *binanceService.Service {
	if srv.binanceService == nil {
		srv.binanceService = binanceService.New(
			srv.AppBinanceClient(ctx),
			srv.AppStockData(ctx).PriceByTicker,
		)
	}

	return srv.binanceService
}

func (srv *serviceProvider) AppHTTPHandler(ctx context.Context) *api.Server {
	if srv.appServer == nil {
		srv.appServer = api.New(
			srv.AppBinanceService(ctx),
		)
	}
	return srv.appServer
}

func (srv *serviceProvider) AppByBitService(ctx context.Context) *bybitService.Service {
	if srv.byBitService == nil {
		srv.byBitService = bybitService.New(
			srv.wsByBitClient,
			srv.AppStockData(ctx).PriceByTicker,
		)
	}

	return srv.byBitService
}

func (srv *serviceProvider) AppWSHandler(ctx context.Context) *wsserver.Server {
	if srv.wsServer == nil {
		srv.wsServer = wsserver.New(
			srv.AppByBitService(ctx),
		)
	}

	return srv.wsServer
}
