package main

import (
	"os"

	"github.com/opendwellers/jujubot/pkg/bot"
	"github.com/opendwellers/jujubot/pkg/config"
	"go.uber.org/zap"
)

var logger *zap.Logger

func initLogger() {
	if os.Getenv("IS_DEBUG") == "true" {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	zap.ReplaceGlobals(logger)
}

// Documentation for the Go driver can be found
// at https://godoc.org/github.com/mattermost/platform/model#Client
func main() {
	// Initialize logger
	initLogger()
	defer logger.Sync()

	// Load configuration
	logger.Sugar().Info("Loading configuration")
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Sugar().Fatal("Failed to load configuration", zap.Error(err))
	}
	logger.Sugar().Info("Configuration loaded")

	// Create and start the bot
	b, err := bot.New(cfg)
	if err != nil {
		logger.Sugar().Fatal("Failed to create bot", zap.Error(err))
	}

	if err := b.Start(); err != nil {
		logger.Sugar().Fatal("Failed to start bot", zap.Error(err))
	}

	// Block forever
	select {}
}
