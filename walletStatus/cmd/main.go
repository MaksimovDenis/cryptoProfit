package main

import (
	"context"
	"os"
	application "walletStatus/internal/app"
	"walletStatus/internal/infra/logger"
)

func main() {
	ctx := context.Background()

	app, err := application.NewApp(ctx, os.Getenv("CONFIG_FILE"))
	if err != nil {
		logger.Fatalf(ctx, "failed to init app: %s", err)
	}

	err = app.Run(ctx)
	if err != nil {
		logger.Fatalf(ctx, "failed to run app: %s", err)
	}
}
