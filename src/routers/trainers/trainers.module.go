package trainers

import (
	"context"
	"errors"
	"math/big"
	"net/http"
	"strings"
	"sync"
	trainers_contract "unrealDestiny/dataAPI/src/routers/trainers/contract"
	"unrealDestiny/dataAPI/src/utils/config"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

func (router *TrainersRouter) getStaticTrainers(c *gin.Context) {
	staticTrainersCollection := router.router.MainDatabase.Collection(COLLECTION_STATIC_TRAINERS)

	cursor, err := staticTrainersCollection.Find(context.TODO(), bson.D{})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
		return
	}

	var results []TrainerStatic

	if err = cursor.All(context.TODO(), &results); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
		return
	}

	c.IndentedJSON(http.StatusOK, results)
}

func (router *TrainersRouter) getUserTrainers(c *gin.Context) {
	walletAddress := c.Param("address")

	staticTrainersCollection := router.router.MainDatabase.Collection(COLLECTION_USER_TRAINERS)

	cursor, err := staticTrainersCollection.Find(context.TODO(), bson.M{"wallet": walletAddress})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
		return
	}

	var results []UserTrainer

	if err = cursor.All(context.TODO(), &results); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
		return
	}

	c.IndentedJSON(http.StatusOK, results)
}

func (router *TrainersRouter) getStaticTrainer(c *gin.Context) {
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

func (router *TrainersRouter) getStaticTrainerDataByModel(model int) (TrainerStatic, error) {
	var trainer TrainerStatic

	staticTrainersCollection := router.router.MainDatabase.Collection(COLLECTION_STATIC_TRAINERS)

	result := staticTrainersCollection.FindOne(context.TODO(), bson.M{"model": model})

	if result.Err() == mongo.ErrNoDocuments {
		return trainer, errors.New("invalid search")
	}

	result.Decode(&trainer)

	return trainer, nil
}

func (router *TrainersRouter) getUserTrainerByIndex(index int) (UserTrainer, error) {
	var trainer UserTrainer

	userTrainersCollection := router.router.MainDatabase.Collection(COLLECTION_USER_TRAINERS)

	result := userTrainersCollection.FindOne(context.TODO(), bson.M{"index": index})

	if result.Err() == mongo.ErrNoDocuments {
		return trainer, result.Err()
	}

	result.Decode(&trainer)

	return trainer, nil
}

// NOTE - addNewUserTrainer(mintingEvent)
// It will generate a new trainer and insert it on the database like document
// It should be execute when the contract launches a new MintTrainer event
func (router *TrainersRouter) addNewUserTrainer(mintingEvent TrainerMinting) {
	var userTrainer UserTrainer

	staticTrainer, err := router.getStaticTrainerDataByModel(int(mintingEvent.Model))

	if err != nil {
		router.router.ServerConfig.LOGGER.Error("Error getting te static trainer data")
		return
	}

	userTrainer.Type = staticTrainer.Type
	userTrainer.Model = staticTrainer.Model
	userTrainer.Experience = 0
	userTrainer.Level = 1
	userTrainer.Health = staticTrainer.Health
	userTrainer.Speed = staticTrainer.Speed
	userTrainer.Energy = staticTrainer.Energy
	userTrainer.Defense = staticTrainer.Defense
	userTrainer.Attack = staticTrainer.Attack
	userTrainer.Builder = staticTrainer.Builder
	userTrainer.Name = staticTrainer.Name
	userTrainer.Network = staticTrainer.Network
	userTrainer.Wallet = mintingEvent.To.String()
	userTrainer.Index = int32(mintingEvent.Token.Int64())

	userTrainersColection := router.router.MainDatabase.Collection(COLLECTION_USER_TRAINERS)

	_, err = userTrainersColection.InsertOne(context.TODO(), userTrainer)

	if err != nil {
		router.router.ServerConfig.LOGGER.Error("Error inserting the new trainer")
		return
	}

	router.router.ServerConfig.LOGGER.Info("New inserted user trainer (" + userTrainer.Wallet + ")")
}

// NOTE - moveTrainerFromOwner(transferEvent)
// Will works to change the owner of the traner when someone executes a transaction on chain
// It is neccessary to change the owner on the offchain interfaces and systems
func (router *TrainersRouter) moveTrainerFromOwner(transferEvent TrainerTransfer) {
	userTrainersCollection := router.router.MainDatabase.Collection(COLLECTION_USER_TRAINERS)

	searchedUserTrainer, err := router.getUserTrainerByIndex(int(transferEvent.Token.Int64()))

	if err != nil {
		if err == mongo.ErrNoDocuments {
			router.router.ServerConfig.LOGGER.Info("Not offchain trainer detected, generating new one offchain")

			instance, err := trainers_contract.NewTrainers(common.HexToAddress(CONTRACT_TRAINERS_ERC721), router.router.ETHCLient)

			if err != nil {
				router.router.ServerConfig.LOGGER.Error("Error creating the trainers contract instance")
			}

			model, err := instance.TokenModel(&bind.CallOpts{}, transferEvent.Token)

			if err != nil {
				router.router.ServerConfig.LOGGER.Error("Error Searching the trainer model on the contract")
			}

			router.addNewUserTrainer(TrainerMinting{Model: uint16(model.Int64()), Token: transferEvent.Token, To: transferEvent.To})

		} else {
			router.router.ServerConfig.LOGGER.Error("Invalid searched trainer")
		}
		return
	}

	if searchedUserTrainer.Wallet != string(transferEvent.From.String()) {
		router.router.ServerConfig.LOGGER.Error("Invalid trainer owner")
		return
	}

	result, err := userTrainersCollection.UpdateOne(context.TODO(), bson.M{"index": int(transferEvent.Token.Int64())}, bson.M{"$set": bson.M{"wallet": transferEvent.To.String()}})

	if err != nil {
		router.router.ServerConfig.LOGGER.Error("Error updating the trainer owner")
		return
	}

	if result.ModifiedCount > 0 {
		router.router.ServerConfig.LOGGER.Info("Updated Trainer owner")
	}
}

func (router *TrainersRouter) initChainListeners() {
	router.router.ServerConfig.LOGGER.Info("Starting Trainers OnChain Listeners")

	logs, sub, err := SubcribeToTransfers(router.router.ETHCLient)

	if err != nil {
		router.router.ServerConfig.LOGGER.Fatal("Trainers Query subscription error")
		return
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(trainers_contract.TrainersABI)))

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
				if IsTransfer(vLog.Topics[0]) {
					var transferEvent TrainerTransfer

					transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
					transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())
					transferEvent.Token = vLog.Topics[3].Big()

					router.router.ServerConfig.LOGGER.Info("Detected Trainer Transfer")

					if transferEvent.From.String() != "0x0000000000000000000000000000000000000000" {
						router.moveTrainerFromOwner(transferEvent)
					}
				} else if IsNewMint(vLog.Topics[0]) {
					var mintEvent TrainerMinting
					var ok bool

					mintTrainerInterface, err := contractAbi.Unpack("MintTrainer", vLog.Data)

					if err != nil {
						router.router.ServerConfig.LOGGER.Error("Error unpacking mint trainer data")
					}

					mintEvent.Model, ok = mintTrainerInterface[0].(uint16)

					if !ok {
						router.router.ServerConfig.LOGGER.Error("Error parsing the trainer minting event data (Model)")
					}

					mintEvent.Token, ok = mintTrainerInterface[1].(*big.Int)

					if !ok {
						router.router.ServerConfig.LOGGER.Error("Error parsing the trainer minting event data (Token)")
					}

					mintEvent.To, ok = mintTrainerInterface[2].(common.Address)

					if !ok {
						router.router.ServerConfig.LOGGER.Error("Error parsing the trainer minting event data (To)")
					}

					router.addNewUserTrainer(mintEvent)
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
	router.router.ParsedGet("/static", router.getStaticTrainers)
	router.router.ParsedGet("/static/:id", router.getStaticTrainer)
	router.router.ParsedGet("/user/:address", router.getUserTrainers)
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
