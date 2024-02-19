package main

import (
	"log"
	"unrealDestiny/dataAPI/src/controller"
	"unrealDestiny/dataAPI/src/utils/config"
	"unrealDestiny/dataAPI/src/utils/env"
	"unrealDestiny/dataAPI/src/utils/logger"

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

	controller.CreateReaderController(&config, router)

	if err := router.Run("localhost:" + config.PORT); err != nil {
		log.Fatal(err)
	}
}
