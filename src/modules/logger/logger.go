package logger

import (
	"log"
	"unrealDestiny/contractsReader/src/modules/config"

	"go.uber.org/zap"
)

func StartLogger(serverConfig *config.ServerConfig) bool {
	zLogger, zErr := zap.NewProduction()

	if zErr != nil {
		log.Fatal("Logger error")
		return false
	}

	defer zLogger.Sync()

	serverConfig.LOGGER = zLogger
	serverConfig.LOGGER.Info("Start logger")

	return true
}
