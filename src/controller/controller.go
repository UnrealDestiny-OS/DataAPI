package controller

import (
	"unrealDestiny/dataAPI/src/routers/users"
	"unrealDestiny/dataAPI/src/utils/config"

	"github.com/gin-gonic/gin"
)

type RoutersConfig struct {
	Users *users.UsersRouter
}

func (config *RoutersConfig) InitAllRoutes() {
	config.Users.CreateRoutes()
}

// NOTE - CreateReaderController(*ServerConfig, *ginEngine)
// Creates all the routers on the application, then manage it to saolve all the gin routes
func CreateReaderController(serverConfig *config.ServerConfig, router *gin.Engine) {
	routers := RoutersConfig{
		Users: users.CreateUsersRouter(serverConfig, router),
	}

	routers.InitAllRoutes()
}
