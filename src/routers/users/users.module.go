package users

import (
	"unrealDestiny/dataAPI/src/utils/config"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersRouter struct {
	router config.Router
}

// SECTION - Internal methods

// NOTE - ParsedGet
// Modify the initial Get function to add the router initial Path, in this case /<InitialRouterPath>/<NewRoutePath>
func (router *UsersRouter) ParsedGet(path string, callback func(*gin.Context)) {
	router.router.MainRouter.GET(router.router.Path+path, callback)
}

// SECTION - REST API
// Rest API methods

func GetAllUsers(context *gin.Context) {

}

func GetAllHolders(context *gin.Context) {

}

// NOTE - Upload all holders
// One time method, it will include all validated holders on the database
func UploadAllHolders(context *gin.Context) {

}

// SECTION - Router Main methods
// All the methods related to the initialization or configuration
// Normally this methods will be called from another core modules

func (router *UsersRouter) CreateRoutes() error {
	router.ParsedGet("/all", GetAllUsers)
	return nil
}

func (router *UsersRouter) Init(serverConfig *config.ServerConfig, mainRouter *gin.Engine, database *mongo.Database) {
	router.router.ServerConfig = serverConfig
	router.router.MainRouter = mainRouter
	router.router.MainDatabase = database
	router.router.Name = "users"
	router.router.Path = "/users"
	router.router.ServerConfig.LOGGER.Info("Starting Users router on " + router.router.Path)
}

func CreateUsersRouter(serverConfig *config.ServerConfig, router *gin.Engine, database *mongo.Database) *UsersRouter {
	users := new(UsersRouter)
	users.Init(serverConfig, router, database)
	return users
}
