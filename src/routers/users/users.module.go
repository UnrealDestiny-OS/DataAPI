package users

import (
	"unrealDestiny/dataAPI/src/utils/config"

	"github.com/gin-gonic/gin"
)

type UsersRouter struct {
	router config.Router
}

func (router *UsersRouter) CreateRoutes() error {
	return nil
}

func (router *UsersRouter) Init(serverConfig *config.ServerConfig, mainRouter *gin.Engine) {
	router.router.ServerConfig = serverConfig
	router.router.MainRouter = mainRouter
	router.router.Name = "users"
	router.router.Path = "/users/"
	router.router.ServerConfig.LOGGER.Info("Starting Users router on " + router.router.Path)
}

func CreateUsersRouter(serverConfig *config.ServerConfig, router *gin.Engine) *UsersRouter {
	users := new(UsersRouter)
	users.Init(serverConfig, router)
	return users
}
