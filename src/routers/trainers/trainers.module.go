package trainers

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"unrealDestiny/dataAPI/src/utils/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (router *TrainersRouter) GetStaticTrainer(c *gin.Context) {
	trainer, err := router.getStaticTrainerData(c.Param("id"))

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, trainer)
}

func (router *TrainersRouter) getStaticTrainerData(id string) (TrainerStatic, error) {
	var trainer TrainerStatic

	staticTrainersCollection := router.router.MainDatabase.Collection(COLLECTION_STATIC_TRAINERS)

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return trainer, errors.New("invalid trainer ID")
	}

	result := staticTrainersCollection.FindOne(context.TODO(), bson.M{"_id": objectId})

	if result.Err() == mongo.ErrNoDocuments {
		return trainer, errors.New("invalid search")
	}

	result.Decode(&trainer)

	return trainer, nil
}

func (router *TrainersRouter) addNewUserTrainer(transferEvent TrainerTransfer) {

}

// NOTE - moveTrainerFromOwner(transferEvent)
// Will works to change the owner of the traner when someone executes a transaction on chain
// It is neccessary to change the owner on the offchain interfaces and systems
func (router *TrainersRouter) moveTrainerFromOwner(transferEvent TrainerTransfer) {
	_ = router.router.MainDatabase.Collection(COLLECTION_USER_TRAINERS)
}

func (router *TrainersRouter) initChainListeners() {
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
						router.addNewUserTrainer(transferEvent)
					} else {
						router.moveTrainerFromOwner(transferEvent)
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
	router.router.ParsedGet("/static", router.GetStaticTrainers)
	router.router.ParsedGet("/static/:id", router.GetStaticTrainer)
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
	router.initChainListeners()
}

func CreateRouter(serverConfig *config.ServerConfig, router *gin.Engine, database *mongo.Database, client *ethclient.Client) *TrainersRouter {
	trainers := new(TrainersRouter)
	trainers.Init(serverConfig, router, database)
	trainers.InitETH(client)
	return trainers
}
