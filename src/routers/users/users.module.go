package users

import (
	"unrealDestiny/dataAPI/src/utils/config"
)

type UsersRouter struct {
	router config.Router
}

func (router *UsersRouter) CreateRoutes() error {
	return nil
}

func (router *UsersRouter) Init(serverConfig *config.ServerConfig) {
	router.router.ServerConfig = serverConfig
	router.router.Name = "users"
	router.router.Path = "/users/"
	router.router.ServerConfig.LOGGER.Info("Starting Users router on " + router.router.Path)
}

func CreateUsersRouter(serverConfig *config.ServerConfig) *UsersRouter {
	router := new(UsersRouter)
	router.Init(serverConfig)
	return router
}