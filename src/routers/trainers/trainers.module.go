package trainers

import (
	"context"
	"net/http"
	"sync"
	"unrealDestiny/dataAPI/src/utils/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrainersRouter struct {
	router config.Router
}

// SECTION - REST API
// Rest API methods

func (router *TrainersRouter) GetStaticTrainers(c *gin.Context) {
	staticTrainersCollection := router.router.MainDatabase.Collection(COLLECTION_STATIC_TRAINERS)

	cursor, err := staticTrainersCollection.Find(context.TODO(), bson.D{})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
	}

	var results []TrainerStatic

	if err = cursor.All(context.TODO(), &results); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
	}

	c.IndentedJSON(http.StatusOK, results)
}

func (router *TrainersRouter) AddNewTrainer(c *gin.Context) {

}

func (router *TrainersRouter) InitChainListeners() {
	router.router.ServerConfig.LOGGER.Info("Starting Trainers OnChain Listeners")

	logs, sub, err := SubcribeToTransfers(router.router.ETHCLient)

	if err != nil {
		router.router.ServerConfig.LOGGER.Fatal("Trainers Query subscription error")
	}

	var wg sync.WaitGroup
	var i int

	go func(i int, wg *sync.WaitGroup) {
		defer wg.Done() // when work is done, declare termination

		for {
			select {
			case <-sub.Err():
				router.router.ServerConfig.LOGGER.Fatal("Error reading the subscription logs")
			case vLog := <-logs:
				if IsTransfer(vLog.Topics[0]) {
					var transferEvent TrainerTransfer

					transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
					transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
					transferEvent.Token = vLog.Topics[3].Big()

					router.router.ServerConfig.LOGGER.Info("Detected Trainer Transfer")

					if transferEvent.From.String() == "0x0000000000000000000000000000000000000000" {
						//New Trainer

					} else {
						// Old trainer new Owner
					}
				}
			}
		}
	}(i, &wg)

	wg.Wait()

}

// SECTION - Router Main methods
// All the methods related to the initialization or configuration
// Normally this methods will be called from another core modules

func (router *TrainersRouter) CreateRoutes() error {
	router.router.ParsedGet("/staticTrainers", router.GetStaticTrainers)
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
	router.InitChainListeners()
}

func CreateRouter(serverConfig *config.ServerConfig, router *gin.Engine, database *mongo.Database, client *ethclient.Client) *TrainersRouter {
	trainers := new(TrainersRouter)
	trainers.Init(serverConfig, router, database)
	trainers.InitETH(client)
	return trainers
}
