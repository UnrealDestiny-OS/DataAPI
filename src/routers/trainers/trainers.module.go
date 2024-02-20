package trainers

import (
	"unrealDestiny/dataAPI/src/utils/config"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrainersRouter struct {
	router config.Router
}

// SECTION - REST API
// Rest API methods

func (router *TrainersRouter) GetStaticTrainers(context *gin.Context) {

}

// SECTION - Router Main methods
// All the methods related to the initialization or configuration
// Normally this methods will be called from another core modules

func (router *TrainersRouter) CreateRoutes() error {
	return nil
}

func (router *TrainersRouter) Init(serverConfig *config.ServerConfig, mainRouter *gin.Engine, database *mongo.Database) {
	router.router.ServerConfig = serverConfig
	router.router.MainRouter = mainRouter
	router.router.MainDatabase = database
	router.router.Name = "trainers"
	router.router.Path = "/trainers"
	router.router.ServerConfig.LOGGER.Info("Starting Trainers router on " + router.router.Path)
}

func (router *TrainersRouter) InitETH(client *ethclient.Client) {
	router.router.ETHCLient = client
	router.router.ServerConfig.LOGGER.Info("Starting Trainers ETH CLient on " + router.router.Path)
}

func CreateRouter(serverConfig *config.ServerConfig, router *gin.Engine, database *mongo.Database, client *ethclient.Client) *TrainersRouter {
	trainers := new(TrainersRouter)
	trainers.Init(serverConfig, router, database)
	trainers.InitETH(client)
	return trainers
}
