package idle

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unrealDestiny/dataAPI/src/routers/users"
	"unrealDestiny/dataAPI/src/utils/config"
	"unrealDestiny/dataAPI/src/utils/contracts"
	"unrealDestiny/dataAPI/src/utils/data"

	"github.com/beeker1121/goque"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/storyicon/sigverify"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IdleExecutorStatus struct {
	available  bool
	privateKey *ecdsa.PrivateKey
	address    common.Address
}

type IdleRouter struct {
	router               config.Router
	deployments          *contracts.Deployments
	executors            []IdleExecutorStatus
	transactionsQueues   []*goque.Queue
	transactionWaitGroup sync.WaitGroup
}

type IdleTransactionQueueElement struct {
	Type         string
	Contract     common.Address
	Data         []byte
	GasLimit     uint64
	GenerationID string
}

// SECTION - Database controllers

func (router *IdleRouter) addNewInjectionLog(injection BalanceInjection, injectionType string) {
	var injectionLog IdleInjectionLog

	userTrainersColection := router.router.MainDatabase.Collection(COLLECTION_IDLE_LOGS_INJECTIONS)

	injectionLog.Amount = injection.Amount.String()
	injectionLog.Wallet = injection.Wallet.String()
	injectionLog.Date = injection.Date.String()
	injectionLog.Type = injectionType
	injectionLog.Network = router.router.ServerConfig.ACTIVE_CHAIN_ID

	_, err := userTrainersColection.InsertOne(context.TODO(), injectionLog)

	if err != nil {
		router.router.ServerConfig.LOGGER.Error("Error inserting the injection log")
		return
	}

	router.router.ServerConfig.LOGGER.Info("New inserted injection (" + injectionLog.Wallet + " - " + injectionLog.Amount + ")")
}

func (router *IdleRouter) updatePlayerInfo(amount string, wallet common.Address) {
	profilesCollection := router.router.MainDatabase.Collection(users.COLLECTION_USER_PROFILES)

	filter := bson.M{"wallet": wallet.String()}
	update := bson.M{"$set": bson.M{"FEE": amount}}

	router.router.ServerConfig.LOGGER.Info("Updating player info (" + wallet.String() + ").")

	_, err := profilesCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		router.router.ServerConfig.LOGGER.Error("Error updating the user balance")
		return
	}
}

// SECTION - Modifiers and executors

func (router *IdleRouter) releaseExecutor(address common.Address) {
	for i := 0; i < len(router.executors); i++ {
		if router.executors[i].address == address {
			router.executors[i].available = true
		}
	}
}

func (router *IdleRouter) getRandomQueue() *goque.Queue {
	min := 0
	max := len(router.transactionsQueues) - 1
	return router.transactionsQueues[rand.Intn(max-min)+min]
}

func (router *IdleRouter) updateTransactionLog(generationID string, hash string, errorString string) {
	transactionsLoggerCollection := router.router.MainDatabase.Collection(COLLECTION_IDLE_EXECUTION_LOG)

	result, err := transactionsLoggerCollection.UpdateOne(context.Background(), bson.M{"generationID": generationID}, bson.M{"$set": bson.M{"hash": hash, "error": errorString}})

	if err != nil {
		router.router.ServerConfig.LOGGER.Error("Error updating the transaction log")
		return
	}

	if result.ModifiedCount > 0 {
		router.router.ServerConfig.LOGGER.Info("Updated transaction log")
	}
}

func (router *IdleRouter) processTransaction(contract common.Address, data []byte, gasLimit uint64, generationID string, executor *IdleExecutorStatus) {
	if executor == nil {
		router.finishTransaction(executor.address)
		router.updateTransactionLog(generationID, "0x", IDLE_NOT_AVAILABLE_EXECUTOR)
		router.router.ServerConfig.LOGGER.Error(IDLE_NOT_AVAILABLE_EXECUTOR)
		return
	} else {
		executor.available = false
		router.router.ServerConfig.LOGGER.Info("Starting the transaction execution (" + generationID + ") using " + executor.address.String())
	}

	nonce, err := router.router.ETHCLient.PendingNonceAt(context.Background(), executor.address)

	if err != nil {
		router.finishTransaction(executor.address)
		router.updateTransactionLog(generationID, "0x", IDLE_ERROR_NONCE)
		router.router.ServerConfig.LOGGER.Error(IDLE_ERROR_NONCE)
		return
	}

	gasPrice, err := router.router.ETHCLient.SuggestGasPrice(context.Background())

	if err != nil {
		router.finishTransaction(executor.address)
		router.updateTransactionLog(generationID, "0x", IDLE_ERROR_GAS_PRICE)
		router.router.ServerConfig.LOGGER.Error(IDLE_ERROR_GAS_PRICE)
		return
	}

	tx := types.NewTransaction(nonce, contract, big.NewInt(0), gasLimit, gasPrice, data)

	chainID, err := router.router.ETHCLient.NetworkID(context.Background())

	if err != nil {
		router.finishTransaction(executor.address)
		router.updateTransactionLog(generationID, "0x", IDLE_ERROR_CHAIN_ID)
		router.router.ServerConfig.LOGGER.Error(IDLE_ERROR_CHAIN_ID)
		return
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), executor.privateKey)

	if err != nil {
		router.finishTransaction(executor.address)
		router.updateTransactionLog(generationID, tx.Hash().String(), IDLE_ERROR_SIGNED_TX)
		router.router.ServerConfig.LOGGER.Error(IDLE_ERROR_SIGNED_TX)
		return
	}

	err = router.router.ETHCLient.SendTransaction(context.Background(), signedTx)

	if err != nil {
		router.finishTransaction(executor.address)
		router.updateTransactionLog(generationID, tx.Hash().String(), IDLE_SENDING_EXECUTOR_TX)
		router.router.ServerConfig.LOGGER.Error(IDLE_SENDING_EXECUTOR_TX)
		return
	}

	receipt, err := bind.WaitMined(context.Background(), router.router.ETHCLient, signedTx)

	if err != nil {
		router.finishTransaction(executor.address)
		router.updateTransactionLog(generationID, tx.Hash().String(), IDLE_WAITING_FOR_MINED)
		router.router.ServerConfig.LOGGER.Error(IDLE_WAITING_FOR_MINED)
		return
	}

	if receipt.Status == 0 {
		router.finishTransaction(executor.address)
		router.updateTransactionLog(generationID, tx.Hash().String(), IDLE_NOT_MINED_TRANSACTION)
		router.router.ServerConfig.LOGGER.Error(IDLE_NOT_MINED_TRANSACTION)
		return
	}

	router.updateTransactionLog(generationID, tx.Hash().String(), "0x")
	router.finishTransaction(executor.address)
}

func (router *IdleRouter) finishTransaction(address common.Address) {
	router.releaseExecutor(address)
	router.transactionWaitGroup.Done()
}

func (router *IdleRouter) validateWalletAuth(senderWallet string, creationChain int, sign string) error {
	valid, err := sigverify.VerifyEllipticCurveHexSignatureEx(
		common.HexToAddress(senderWallet),
		[]byte(GenerateTransactionExecutionSign(senderWallet, creationChain)),
		sign,
	)

	if err != nil {
		return errors.New(IDLE_ERROR_VALIDATING_TX_SIGN)
	}

	if !valid {
		return errors.New(IDLE_INVALID_SIGN)
	}

	return nil
}

func (router *IdleRouter) validateUserFee(senderWallet string) error {
	profilesCollection := router.router.MainDatabase.Collection(users.COLLECTION_USER_PROFILES)

	var result users.UserProfile
	err := profilesCollection.FindOne(context.Background(), bson.M{"wallet": senderWallet}).Decode(&result)

	if err != nil {
		return errors.New(IDLE_ERROR_SEARCHING_USER_INFO)
	}

	wei := new(big.Int)
	wei.SetString(result.FEE, 10)
	userBalance := data.ToDecimal(wei, 18)
	needBalance, err := decimal.NewFromString("0.05")

	if err != nil {
		return errors.New(IDLE_DECIMAL_FORMAT_PARSING_ERROR)
	}

	if userBalance.LessThan(needBalance) {
		return errors.New(IDLE_NOT_ENOUGHT_FEE)
	}

	return nil
}

// SECTION - REST API
// Rest API methods

func (router *IdleRouter) executeTrainerJoinRequest(c *gin.Context) {
	var trainerJoin APIExecuteTrainerJoin

	err := c.BindJSON(&trainerJoin)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": IDLE_ERROR_TRAINER_JOIN_PARSING})
		return
	}

	err = router.validateWalletAuth(trainerJoin.Wallet, router.router.ServerConfig.ACTIVE_CHAIN_ID, trainerJoin.WalletAuth)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	err = router.validateUserFee(trainerJoin.Wallet)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	contractAddress := common.HexToAddress(router.deployments.TrainersIDLE.Address)
	trainerJoinTxData := GetTrainerJoinTxData(trainerJoin.Wallet, trainerJoin.Trainer)

	txGenerationID, err := router.initializeDatabaseTransaction()

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	_, err = router.getRandomQueue().EnqueueObject(IdleTransactionQueueElement{Type: "TRAINER_JOIN", Contract: contractAddress, Data: trainerJoinTxData, GasLimit: uint64(110000), GenerationID: *txGenerationID})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": IDLE_ENQUEUE_ERROR})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"error": false, "generationID": *txGenerationID})
}

func (router *IdleRouter) executeCollectTransactionPointsRequest(c *gin.Context) {
	var collectTransactionPoints APIExecuteCollectTransactionPoints

	err := c.BindJSON(&collectTransactionPoints)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": IDLE_ERROR_REQUEST_PARSING})
		return
	}

	err = router.validateWalletAuth(collectTransactionPoints.Wallet, router.router.ServerConfig.ACTIVE_CHAIN_ID, collectTransactionPoints.WalletAuth)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	err = router.validateUserFee(collectTransactionPoints.Wallet)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	contractAddress := common.HexToAddress(router.deployments.TrainersIDLE.Address)
	trainerJoinTxData := GetCollectTransactionPointsData(collectTransactionPoints.Wallet, collectTransactionPoints.Trainer)

	txGenerationID, err := router.initializeDatabaseTransaction()

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	_, err = router.getRandomQueue().EnqueueObject(IdleTransactionQueueElement{Type: "COLLECT_TRANSACTION_POINTS", Contract: contractAddress, Data: trainerJoinTxData, GenerationID: *txGenerationID, GasLimit: uint64(110000)})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": IDLE_ENQUEUE_ERROR})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"error": false, "generationID": *txGenerationID})
}

func (router *IdleRouter) executeCollectIdlePointsRequest(c *gin.Context) {
	var collectTransactionPoints APIExecuteCollectTransactionPoints

	err := c.BindJSON(&collectTransactionPoints)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": IDLE_ERROR_REQUEST_PARSING})
		return
	}

	err = router.validateWalletAuth(collectTransactionPoints.Wallet, router.router.ServerConfig.ACTIVE_CHAIN_ID, collectTransactionPoints.WalletAuth)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	err = router.validateUserFee(collectTransactionPoints.Wallet)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	contractAddress := common.HexToAddress(router.deployments.TrainersIDLE.Address)
	trainerJoinTxData := GetCollectionIdlePointsData(collectTransactionPoints.Wallet, collectTransactionPoints.Trainer)

	txGenerationID, err := router.initializeDatabaseTransaction()

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	_, err = router.getRandomQueue().EnqueueObject(IdleTransactionQueueElement{Type: "COLLECT_TRANSACTION_POINTS", Contract: contractAddress, Data: trainerJoinTxData, GenerationID: *txGenerationID, GasLimit: uint64(110000)})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": IDLE_ENQUEUE_ERROR})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"error": false, "generationID": *txGenerationID})
}

func (router *IdleRouter) executeGetImprovementRequest(c *gin.Context) {
	var buyImprovement APIExecuteImprovementBuy

	err := c.BindJSON(&buyImprovement)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": IDLE_ERROR_REQUEST_PARSING})
		return
	}

	err = router.validateWalletAuth(buyImprovement.Wallet, router.router.ServerConfig.ACTIVE_CHAIN_ID, buyImprovement.WalletAuth)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	err = router.validateUserFee(buyImprovement.Wallet)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	contractAddress := common.HexToAddress(router.deployments.TrainersIDLE.Address)
	trainerJoinTxData := GetBuyImprovementData(buyImprovement.Wallet, buyImprovement.Trainer, buyImprovement.Improvement)

	txGenerationID, err := router.initializeDatabaseTransaction()

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	_, err = router.getRandomQueue().EnqueueObject(IdleTransactionQueueElement{Type: "BUY_IMPROVEMENT", Contract: contractAddress, Data: trainerJoinTxData, GenerationID: *txGenerationID, GasLimit: uint64(110000)})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": IDLE_ENQUEUE_ERROR})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"error": false, "generationID": *txGenerationID})
}

// SECTION - Collection Idle points

func (router *IdleRouter) initializeDatabaseTransaction() (*string, error) {
	var transactionLog TransactionExecutionLog
	txGenerationID := "tx-" + strconv.Itoa(int(time.Now().UnixNano())) + "-" + strconv.Itoa(int(rand.Intn(10000000)))
	transactionsLoggerCollection := router.router.MainDatabase.Collection(COLLECTION_IDLE_EXECUTION_LOG)
	_, err := transactionsLoggerCollection.InsertOne(context.TODO(), transactionLog)

	if err != nil {
		return nil, errors.New(IDLE_DATABASE_INSERT_ERROR)
	}

	return &txGenerationID, nil
}

// SECTION - Onchain Listeners

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
				if IsBalanceInjection(vLog.Topics[0]) || IsBalanceWithdraw(vLog.Topics[0]) {
					var injectionEvent BalanceInjection
					var injectionType = "INJECTION"
					var eventType = "InjectFeeBalance"

					if IsBalanceWithdraw(vLog.Topics[0]) {
						injectionType = "WITHDRAW"
						eventType = "WithdrawBalance"
					}

					err := ProcessInjectionEvent(eventType, contractAbi, vLog.Data, &injectionEvent)

					if err != nil {
						router.router.ServerConfig.LOGGER.Error(err.Error())
					}

					router.addNewInjectionLog(injectionEvent, injectionType)
					router.updatePlayerInfo(injectionEvent.Amount.String(), injectionEvent.Wallet)
				} else if IsTakeFees(vLog.Topics[0]) {
					router.router.ServerConfig.LOGGER.Info("Detect new contract event (TAKE FEES) " + vLog.Topics[0].String())
					var takeFeesEvent TakeFees

					err := ProcessTakeFeesEvent(contractAbi, vLog.Data, &takeFeesEvent)

					if err != nil {
						router.router.ServerConfig.LOGGER.Error(err.Error())
					}

					router.updatePlayerInfo(takeFeesEvent.WalletBalance.String(), takeFeesEvent.Wallet)
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
	router.router.ParsedPost("/trainer-join", router.executeTrainerJoinRequest)
	router.router.ParsedPost("/collect-transaction-points", router.executeCollectTransactionPointsRequest)
	router.router.ParsedPost("/collect-idle-points", router.executeCollectIdlePointsRequest)
	router.router.ParsedPost("/get-improvement", router.executeGetImprovementRequest)
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

func (router *IdleRouter) InitExecutors() {
	privateAccounts := router.router.ServerConfig.EXECUTOR_PRIVATE_KEYS

	if len(privateAccounts) > 0 {
		for i := 0; i < len(privateAccounts); i++ {
			var executor IdleExecutorStatus

			executor.privateKey = privateAccounts[i]

			publicKey := executor.privateKey.Public()

			publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

			if !ok {
				log.Fatal("error casting public key to ECDSA")
				break
			}

			executor.address = crypto.PubkeyToAddress(*publicKeyECDSA)
			executor.available = true

			router.router.ServerConfig.LOGGER.Info("Load " + executor.address.String() + " as executor")
			router.executors = append(router.executors, executor)
		}
	} else {
		router.router.ServerConfig.LOGGER.Fatal("Error. Invalid executors.")
	}
}

func (router *IdleRouter) ValidateQueue(index int) {
	for {
		item, err := router.transactionsQueues[index].Peek()

		if err == goque.ErrDBClosed && err == goque.ErrIncompatibleType {
			router.router.ServerConfig.LOGGER.Error("Invalid queue creation")
			return
		} else {
			if item != nil {
				router.transactionWaitGroup.Add(1)
				var transactionElement IdleTransactionQueueElement
				item.ToObject(&transactionElement)
				router.router.ServerConfig.LOGGER.Info("Start a transaction execution using a queue (" + strconv.Itoa(index) + ") - " + transactionElement.Type)
				go router.processTransaction(transactionElement.Contract, transactionElement.Data, transactionElement.GasLimit, transactionElement.GenerationID, &router.executors[index])
				router.transactionWaitGroup.Wait()
				router.router.ServerConfig.LOGGER.Info("Finish transaction execution")
				router.transactionsQueues[index].Dequeue()
			}
		}
	}
}

func (router *IdleRouter) InitTransactionQueues() {
	wd, err := os.Getwd()

	if err != nil {
		router.router.ServerConfig.LOGGER.Fatal("Error. Invalid working dir on queue building.")
	}

	for i := 0; i < len(router.executors); i++ {
		q, err := goque.OpenQueue(wd + "/src/data/queues/" + router.executors[i].address.String())

		if err != nil {
			fmt.Println(err)
			router.router.ServerConfig.LOGGER.Fatal("Error. Invalid transactions queue.")
		}

		router.transactionsQueues = append(router.transactionsQueues, q)

		go router.ValidateQueue(i)
	}

}

func (router *IdleRouter) InitETH(client *ethclient.Client, deployments *contracts.Deployments) {
	router.router.ETHCLient = client
	router.deployments = deployments
	router.router.ServerConfig.LOGGER.Info("Starting Idle ETH CLient on " + router.router.Path)
	router.InitExecutors()
	router.initChainListeners()
	router.InitTransactionQueues()
}

func CreateRouter(serverConfig *config.ServerConfig, router *gin.Engine, database *mongo.Database, client *ethclient.Client, deployments *contracts.Deployments) *IdleRouter {
	trainers := new(IdleRouter)
	trainers.Init(serverConfig, router, database)
	trainers.InitETH(client, deployments)
	return trainers
}
