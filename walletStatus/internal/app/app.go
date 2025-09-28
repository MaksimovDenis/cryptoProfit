package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	bybitws "walletStatus/internal/adapter/bybit_ws"
	roundtripper "walletStatus/internal/adapter/round_tripper"
	"walletStatus/internal/infra/config"
	"walletStatus/internal/infra/logger"

	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
)

const retryCount = 3

type App struct {
	config          *config.Config
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

func NewApp(ctx context.Context, configPath string) (*App, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("config.LoadConfig: %w", err)
	}

	app := &App{config: cfg}

	err = app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	logger.Infof(ctx, "Config loaded successfully")

	return app, nil
}

func (app *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.runHTTPServer(ctx); err != nil {
			logger.Errorf(ctx, "http server error: %v", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.runWSConnectoin(ctx); err != nil {
			logger.Errorf(ctx, "ws connection error: %v", err)
			cancel()
		}
	}()

	gracefulShutdown(ctx, cancel, wg)

	return nil
}

func (app *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		app.initLogger,
		app.initServiceProvider,
		app.initByBitWsClient,
		app.initHTTPBinanceClient,
		app.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *App) initLogger(_ context.Context) error {
	var level zapcore.Level

	switch {
	case app.config.Server.LogLevel == "debug":
		level = zapcore.DebugLevel
	case app.config.Server.LogLevel == "info":
		level = zapcore.InfoLevel
	case app.config.Server.LogLevel == "error":
		level = zapcore.ErrorLevel
	case app.config.Server.LogLevel == "warn":
		level = zapcore.ErrorLevel
	case app.config.Server.LogLevel == "fatal":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.ErrorLevel
	}

	err := logger.Init(level)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) initServiceProvider(_ context.Context) error {
	app.serviceProvider = newServiceProvider()
	app.serviceProvider.config = *app.config
	return nil
}

func (app *App) initHTTPServer(ctx context.Context) error {
	router := app.serviceProvider.AppHTTPHandler(ctx).InitRoutes()

	address := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.Port)

	app.httpServer = &http.Server{
		Addr:              address,
		Handler:           router,
		ReadHeaderTimeout: time.Duration(app.config.Server.Timeout) * time.Second,
	}

	return nil
}

func (app *App) initHTTPBinanceClient(_ context.Context) error {
	transport := http.DefaultTransport

	rateLimiter := rate.NewLimiter(
		rate.Limit(app.config.Binance.Limit),
		app.config.Binance.Burst,
	)

	rateLimitedTransport := roundtripper.NewRateLimitRoundTripper(transport, rateLimiter)

	retryTransport := roundtripper.New(
		rateLimitedTransport,
		retryCount,
		time.Duration(app.config.Server.RetryTimeout)*time.Second,
	)

	app.serviceProvider.httpBinanceClient = &http.Client{
		Transport: retryTransport,
		Timeout:   time.Duration(app.config.Binance.Timeout) * time.Second,
	}

	return nil
}

func (app *App) runHTTPServer(ctx context.Context) error {
	logger.Infof(ctx, "HTTP server is running on %s", app.httpServer.Addr)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.httpServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		logger.Infof(ctx, "HTTP server shutdown initiated")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return app.httpServer.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

func (app *App) initByBitWsClient(ctx context.Context) error {
	var err error

	app.serviceProvider.wsByBitClient, err = bybitws.New(ctx, "wss://stream.bybit.com/v5/public/spot")
	if err != nil {
		return err
	}

	return nil
}

func (app *App) runWSConnectoin(ctx context.Context) error {
	address := fmt.Sprintf("%s:%s", app.config.Server.Host, app.config.Server.WsPort)

	logger.Infof(ctx, "WebSocket connection is running on %s", address)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.serviceProvider.AppWSHandler(ctx).Run(address)
	}()

	select {
	case <-ctx.Done():
		logger.Infof(ctx, "HTTP server shutdown initiated")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return app.httpServer.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}

func gracefulShutdown(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		logger.Infof(ctx, "terminating: context cancelled")
	case <-waitSignal():
		logger.Infof(ctx, "terminating: via signal")
	}

	cancel()
	wg.Wait()
}

func waitSignal() chan os.Signal {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	return sigterm
}
