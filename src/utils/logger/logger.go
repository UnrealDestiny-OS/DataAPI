package logger

import (
	"log"
	"unrealDestiny/contractsReader/src/utils/config"

	"go.uber.org/zap"
)

// NOTE - StartLogger (*ServerConfig)
// Loads the logger using the zap library
// Then set it on the ServerConfig struct
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
