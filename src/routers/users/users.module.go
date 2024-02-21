package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"unrealDestiny/dataAPI/src/utils/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/gin-gonic/gin"
	"github.com/storyicon/sigverify"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersRouter struct {
	router config.Router
}

// SECTION - REST API
// Rest API methods

// NOTE - AddPossibleUser
// GET Request, No Body, No params
func (router *UsersRouter) GetPossibleUsers(c *gin.Context) {
	possibleUsersCollection := router.router.MainDatabase.Collection(COLLECTION_POSSIBLE_USERS)

	var users []UserStaticPossible

	cursor, err := possibleUsersCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
		return
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

// NOTE - AddPossibleUser
// POST Request, Body *PossibleUser
// Insert new possible user when the user reach the web page and connect the wallet to the site
func (router *UsersRouter) AddPossibleUser(c *gin.Context) {
	var user UserStaticPossible
	var searchedUser UserStaticPossible

	err := c.BindJSON(&user)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
		fmt.Print(err)
		return
	}

	possibleUsersCollection := router.router.MainDatabase.Collection(COLLECTION_POSSIBLE_USERS)

	err = possibleUsersCollection.FindOne(context.TODO(), bson.M{"address": user.Address}).Decode(&searchedUser)

	if err != mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": "The user already exists"})
		return
	}

	result, err := possibleUsersCollection.InsertOne(context.TODO(), user)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

// NOTE - GetAllHolders
// GET Request, No Body, No params
// Return all users data on the holders database collection
func (router *UsersRouter) GetAllHolders(c *gin.Context) {
	holdersCollection := router.router.MainDatabase.Collection(COLLECTION_HOLDERS)

	var users []UserStaticHolder

	cursor, err := holdersCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
		return
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

// SECTION - Profiles
// Users profiles

func (router *UsersRouter) createProfileUsingSign(c *gin.Context) {
	var createProfileRequest APICreateUserProfile

	err := c.BindJSON(&createProfileRequest)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
		return
	}

	domain := apitypes.TypedDataDomain{
		Name:              "unrealdestinycom",
		Version:           "1",
		ChainId:           math.NewHexOrDecimal256(createProfileRequest.CreationChain),
		VerifyingContract: "0x0000000000000000000000000000000000000000",
	}

	types := apitypes.Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"Request": {
			{Name: "wallet", Type: "address"},
			{Name: "chain", Type: "uint256"},
			{Name: "username", Type: "string"},
		},
	}

	message := SignCreateUserRequest{
		Wallet:   common.HexToAddress(createProfileRequest.Wallet),
		Username: createProfileRequest.Username,
		Chain:    int(createProfileRequest.CreationChain),
	}

	typedData := apitypes.TypedData{
		Types:       types,
		PrimaryType: "ERC721Order",
		Domain:      domain,
		Message:     message,
	}

	if err := json.Unmarshal([]byte(userData), &typedData); err != nil {
		panic(err)
	}

	fmt.Println(userData)

	valid, err := sigverify.VerifyTypedDataHexSignatureEx(
		common.HexToAddress(createProfileRequest.Wallet),
		typedData,
		createProfileRequest.Sign,
	)

	fmt.Println(valid, err)

}

// SECTION - Router Main methods
// All the methods related to the initialization or configuration
// Normally this methods will be called from another core modules

func (router *UsersRouter) CreateRoutes() error {
	router.router.ParsedGet("/holders", router.GetAllHolders)
	router.router.ParsedPost("/possible", router.AddPossibleUser)
	router.router.ParsedGet("/possible", router.GetPossibleUsers)
	// router.router.ParsedGet("/profile/:wallet", router.GetPossibleUsers)
	router.router.ParsedPost("/profile", router.createProfileUsingSign)
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
