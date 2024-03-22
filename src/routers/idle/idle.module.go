package idle

import (
	"context"
	"math/big"
	"strings"
	"sync"
	"unrealDestiny/dataAPI/src/routers/users"
	"unrealDestiny/dataAPI/src/utils/config"
	"unrealDestiny/dataAPI/src/utils/contracts"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IdleRouter struct {
	router      config.Router
	deployments *contracts.Deployments
}

// SECTION - Database controllers

func (router *IdleRouter) addNewInjectionLog(injection BalanceInjection) {
	var injectionLog IdleInjectionLog

	userTrainersColection := router.router.MainDatabase.Collection(COLLECTION_IDLE_LOGS_INJECTIONS)

	injectionLog.Amount = injection.Amount.String()
	injectionLog.Wallet = injection.Wallet.String()
	injectionLog.Date = injection.Date.String()
	injectionLog.Network = router.router.ServerConfig.ACTIVE_CHAIN_ID

	_, err := userTrainersColection.InsertOne(context.TODO(), injectionLog)

	if err != nil {
		router.router.ServerConfig.LOGGER.Error("Error inserting the injection log")
		return
	}

	router.router.ServerConfig.LOGGER.Info("New inserted injection (" + injectionLog.Wallet + " - " + injectionLog.Amount + ")")
}

func (router *IdleRouter) updatePlayerInfo(injection BalanceInjection, wallet common.Address) {
	profilesCollection := router.router.MainDatabase.Collection(users.COLLECTION_USER_PROFILES)

	filter := bson.M{"wallet": wallet.String()}
	update := bson.M{"$set": bson.M{"FEE": injection.Amount.String()}}

	_, err := profilesCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		router.router.ServerConfig.LOGGER.Error("Error updating the user balance")
		return
	}
}

// SECTION - REST API
// Rest API methods

func (router *IdleRouter) initChainListeners() {
	router.router.ServerConfig.LOGGER.Info("Starting Idle OnChain Listeners using " + router.deployments.TrainersIDLE.Address)

	logs, sub, err := SubscribeToContract(router.router.ETHCLient, router.deployments.TrainersIDLE.Address)

	if err != nil {
		router.router.ServerConfig.LOGGER.Fatal("Idle Query subscription error")
		return
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(*router.deployments.TrainersIDLE.JsonAbi())))

	if err != nil {
		router.router.ServerConfig.LOGGER.Fatal("Error parsing the contract ABI")
		return
	}

	var wg sync.WaitGroup
	var i int

	go func(i int, wg *sync.WaitGroup) {
		defer wg.Done() // when work is done, declare termination

		for {
			select {
			case <-sub.Err():
				router.router.ServerConfig.LOGGER.Error("Error reading the subscription logs")
			case vLog := <-logs:
				if IsBalanceInjection(vLog.Topics[0]) {
					var injectionEvent BalanceInjection
					var ok bool

					injectBalanceInterface, err := contractAbi.Unpack("InjectFeeBalance", vLog.Data)

					if err != nil {
						router.router.ServerConfig.LOGGER.Error("Error unpacking mint trainer data")
						return
					}

					injectionEvent.Wallet, ok = injectBalanceInterface[0].(common.Address)

					if !ok {
						router.router.ServerConfig.LOGGER.Error("Error parsing the injection event data (Wallet)")
						return
					}

					injectionEvent.Amount, ok = injectBalanceInterface[1].(*big.Int)

					if !ok {
						router.router.ServerConfig.LOGGER.Error("Error parsing the injection event data (Amount)")
						return
					}

					injectionEvent.Date, ok = injectBalanceInterface[2].(*big.Int)

					if !ok {
						router.router.ServerConfig.LOGGER.Error("Error parsing the injection event data (Date)")
						return
					}

					router.addNewInjectionLog(injectionEvent)
					router.updatePlayerInfo(injectionEvent, injectionEvent.Wallet)
				}
			}
		}
	}(i, &wg)

	wg.Wait()
}

// SECTION - Router Main methods
// All the methods related to the initialization or configuration
// Normally this methods will be called from another core modules

func (router *IdleRouter) CreateRoutes() error {
	return nil
}

func (router *IdleRouter) Init(serverConfig *config.ServerConfig, mainRouter *gin.Engine, database *mongo.Database) {
	router.router.ServerConfig = serverConfig
	router.router.MainRouter = mainRouter
	router.router.MainDatabase = database
	router.router.Name = "idle"
	router.router.Path = "/idle"
	router.router.ServerConfig.LOGGER.Info("Starting Idle router on " + router.router.Path)
}

func (router *IdleRouter) InitETH(client *ethclient.Client, deployments *contracts.Deployments) {
	router.router.ETHCLient = client
	router.deployments = deployments
	router.router.ServerConfig.LOGGER.Info("Starting Idle ETH CLient on " + router.router.Path)
	router.initChainListeners()
}

func CreateRouter(serverConfig *config.ServerConfig, router *gin.Engine, database *mongo.Database, client *ethclient.Client, deployments *contracts.Deployments) *IdleRouter {
	trainers := new(IdleRouter)
	trainers.Init(serverConfig, router, database)
	trainers.InitETH(client, deployments)
	return trainers
}
