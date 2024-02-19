package controller

import (
	"errors"
	"unrealDestiny/dataAPI/src/routers/users"
	"unrealDestiny/dataAPI/src/utils/config"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var USERS_ROUTER_ERROR string = "Error starting the users router"

type RoutersConfig struct {
	Users *users.UsersRouter
}

func (config *RoutersConfig) InitAllRoutes(serverConfig *config.ServerConfig) error {
	err := config.Users.CreateRoutes()

	if err != nil {
		serverConfig.LOGGER.Fatal(USERS_ROUTER_ERROR)
		return errors.New(USERS_ROUTER_ERROR)
	}

	return nil
}

// NOTE - CreateReaderController(*ServerConfig, *ginEngine)
// Creates all the routers on the application, then manage it to saolve all the gin routes
func CreateReaderController(serverConfig *config.ServerConfig, router *gin.Engine, databaseClient *mongo.Client, database *mongo.Database) error {
	routers := RoutersConfig{
		Users: users.CreateUsersRouter(serverConfig, router, database),
	}

	err := routers.InitAllRoutes(serverConfig)

	if err != nil {
		return err
	}

	return nil
}
