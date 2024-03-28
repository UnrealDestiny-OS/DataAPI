package controller

import (
	"errors"
	"unrealDestiny/dataAPI/src/routers/idle"
	runtime_router "unrealDestiny/dataAPI/src/routers/runtime"
	"unrealDestiny/dataAPI/src/routers/trainers"
	"unrealDestiny/dataAPI/src/routers/users"
	"unrealDestiny/dataAPI/src/utils/config"
	"unrealDestiny/dataAPI/src/utils/contracts"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var USERS_ROUTER_ERROR string = "Error starting the users router"
var TRAINERS_ROUTER_ERROR string = "Error starting the trainers router"
var IDLE_ROUTERS_ERROR string = "Error starting the idle router"
var RUNTIME_ROUTER_ERROR string = "Error starting the runtime router"

type RoutersConfig struct {
	Users    *users.UsersRouter
	Runtime  *runtime_router.RuntimeRouter
	Trainers *trainers.TrainersRouter
	Idle     *idle.IdleRouter
}

func (config *RoutersConfig) InitAllRoutes(serverConfig *config.ServerConfig) error {
	err := config.Users.CreateRoutes()

	if err != nil {
		serverConfig.LOGGER.Fatal(USERS_ROUTER_ERROR)
		return errors.New(USERS_ROUTER_ERROR)
	}

	err = config.Runtime.CreateRoutes()

	if err != nil {
		serverConfig.LOGGER.Fatal(RUNTIME_ROUTER_ERROR)
		return errors.New(RUNTIME_ROUTER_ERROR)
	}

	err = config.Trainers.CreateRoutes()

	if err != nil {
		serverConfig.LOGGER.Fatal(TRAINERS_ROUTER_ERROR)
		return errors.New(TRAINERS_ROUTER_ERROR)
	}

	err = config.Idle.CreateRoutes()

	if err != nil {
		serverConfig.LOGGER.Fatal(IDLE_ROUTERS_ERROR)
		return errors.New(IDLE_ROUTERS_ERROR)
	}

	return nil
}

// NOTE - CreateReaderController(*ServerConfig, *ginEngine)
// Creates all the routers on the application, then manage it to saolve all the gin routes
func CreateReaderController(serverConfig *config.ServerConfig, router *gin.Engine, databaseClient *mongo.Client, database *mongo.Database, client *ethclient.Client) error {
	contractDeployments := contracts.LoadDeploymentsData(serverConfig.USE_PRODUCTION_ADDRESSES)

	routers := RoutersConfig{
		Users:    users.CreateUsersRouter(serverConfig, router, database),
		Runtime:  runtime_router.CreateRouter(serverConfig, router, database),
		Trainers: trainers.CreateRouter(serverConfig, router, database, client, contractDeployments),
		Idle:     idle.CreateRouter(serverConfig, router, database, client, contractDeployments),
	}

	err := routers.InitAllRoutes(serverConfig)

	if err != nil {
		return err
	}

	return nil
}
