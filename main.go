package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"unrealDestiny/dataAPI/src/controller"
	"unrealDestiny/dataAPI/src/utils/config"
	"unrealDestiny/dataAPI/src/utils/database"
	"unrealDestiny/dataAPI/src/utils/env"
	"unrealDestiny/dataAPI/src/utils/logger"
	"unrealDestiny/dataAPI/src/utils/network"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	config := config.ServerConfig{
		CONTEXT: ctx,
	}

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

	// SECTION Database

	databaseClient, database := database.InitDatabase(&config)

	if database == nil {
		log.Fatal("Database error")
		return
	}

	// SECTION ETH

	client, err := ethclient.Dial(config.MTRG_WS_CLIENT)

	if err != nil {
		fmt.Println(err)
		config.LOGGER.Fatal("ETH Client errorr")
	}

	// SECTION Server

	router := gin.Default()

	router.Use(network.CORSMiddleware())

	config.LOGGER.Info("Starting server on localhost:" + config.PORT)

	err = controller.CreateReaderController(&config, router, databaseClient, database, client)

	if err != nil {
		config.LOGGER.Fatal("Routers errorr")
		return
	}

	router.Run("localhost:" + config.PORT)
}
