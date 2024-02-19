package main

import (
	"log"
	"unrealDestiny/contractsReader/src/controller"
	"unrealDestiny/contractsReader/src/utils/config"
	"unrealDestiny/contractsReader/src/utils/env"
	"unrealDestiny/contractsReader/src/utils/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	config := config.ServerConfig{}

	// SECTION ENV variables

	if !env.LoadEnv(&config) {
		log.Fatal("Env error")
		return
	}

	// SECTION Logger

	if !logger.StartLogger(&config) {
		log.Fatal("Logger error")
		return
	}

	// SECTION Server

	router := gin.Default()

	config.LOGGER.Info("Starting server on localhost:" + config.PORT)

	controller.ReaderController(&config, router)

	router.Run("localhost:" + config.PORT)
}
