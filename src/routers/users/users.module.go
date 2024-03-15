package users

import (
	"context"
	"net/http"
	"unrealDestiny/dataAPI/src/utils/config"

	"github.com/ethereum/go-ethereum/common"
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
func (router *UsersRouter) getPossibleUsers(c *gin.Context) {
	possibleUsersCollection := router.router.MainDatabase.Collection(COLLECTION_POSSIBLE_USERS)

	var users []UserStaticPossible

	cursor, err := possibleUsersCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": USERS_ERROR_DATABASE_ERROR})
		return
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": USERS_ERROR_PARSING})
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

// NOTE - AddPossibleUser
// POST Request, Body *PossibleUser
// Insert new possible user when the user reach the web page and connect the wallet to the site
func (router *UsersRouter) addPossibleUser(c *gin.Context) {
	var user UserStaticPossible
	var searchedUser UserStaticPossible

	err := c.BindJSON(&user)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": USERS_ERROR_PARSING})
		return
	}

	possibleUsersCollection := router.router.MainDatabase.Collection(COLLECTION_POSSIBLE_USERS)

	err = possibleUsersCollection.FindOne(context.TODO(), bson.M{"address": user.Address}).Decode(&searchedUser)

	if err != mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": USERS_ERROR_THE_USER_EXISTS})
		return
	}

	result, err := possibleUsersCollection.InsertOne(context.TODO(), user)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": USERS_ERROR_INSERT_DATABASE_ERROR})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

// NOTE - GetAllHolders
// GET Request, No Body, No params
// Return all users data on the holders database collection
func (router *UsersRouter) getAllHolders(c *gin.Context) {
	holdersCollection := router.router.MainDatabase.Collection(COLLECTION_HOLDERS)

	var users []UserStaticHolder

	cursor, err := holdersCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": USERS_ERROR_DATABASE_ERROR})
		return
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": USERS_ERROR_PARSING})
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

// SECTION - Profiles
// Users profiles

// NOTE - createProfileUsingSign
// Create new user profile using the signature verification
// The new user will be created only if the signature is valid compared with the sent wallet address
func (router *UsersRouter) createProfileUsingSign(c *gin.Context) {
	var createProfileRequest APICreateUserProfile

	err := c.BindJSON(&createProfileRequest)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": USERS_ERROR_PARSING})
		return
	}

	valid, err := sigverify.VerifyEllipticCurveHexSignatureEx(
		common.HexToAddress(createProfileRequest.Wallet),
		[]byte(GenerateSignValidationMessage(createProfileRequest.Wallet, int(createProfileRequest.CreationChain))),
		createProfileRequest.Sign,
	)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": USERS_ERROR_ON_SING_VERIFICATION})
		return
	}

	if !valid {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": USERS_ERROR_INVALID_SIGN})
		return
	}

	user := UserProfile{
		Wallet:     createProfileRequest.Wallet,
		Name:       createProfileRequest.Username,
		UDT:        0,
		FUDT:       0,
		Level:      1,
		Experience: 0,
		Chain:      createProfileRequest.CreationChain,
	}

	var searchedUser UserProfile

	profilesCollection := router.router.MainDatabase.Collection(COLLECTION_USER_PROFILES)

	err = profilesCollection.FindOne(context.TODO(), bson.M{"wallet": createProfileRequest.Wallet}).Decode(&searchedUser)

	if err != mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": USERS_ERROR_THE_USER_EXISTS})
		return
	}

	result, err := profilesCollection.InsertOne(context.TODO(), user)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": USERS_ERROR_INSERT_DATABASE_ERROR})
		return
	}

	router.router.ServerConfig.LOGGER.Info("New user profile " + createProfileRequest.Username)
	c.IndentedJSON(http.StatusOK, result)
}

// NOTE - getUserProfile
func (router *UsersRouter) getUserProfile(c *gin.Context) {
	var user UserProfile

	profileSCollection := router.router.MainDatabase.Collection(COLLECTION_USER_PROFILES)

	result := profileSCollection.FindOne(context.TODO(), bson.M{"wallet": c.Param("address")})

	if result.Err() == mongo.ErrNoDocuments {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": USERS_ERROR_THE_USER_NOT_EXISTS})
		return
	}

	result.Decode(&user)

	c.IndentedJSON(http.StatusOK, user)
}

func (router *UsersRouter) validateProfile(c *gin.Context) {
	var validateProfileRequest APIValidateUserProfile
	var user UserProfile

	err := c.BindJSON(&validateProfileRequest)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": USERS_ERROR_PARSING})
		return
	}

	profilesCollection := router.router.MainDatabase.Collection(COLLECTION_USER_PROFILES)

	result := profilesCollection.FindOne(context.TODO(), bson.M{"wallet": c.Param("address")})

	if result.Err() == mongo.ErrNoDocuments {
		user = UserProfile{
			Wallet:     validateProfileRequest.Wallet,
			Name:       "Anonymous ",
			UDT:        0,
			FUDT:       0,
			Level:      1,
			Experience: 0,
			Chain:      validateProfileRequest.CreationChain,
		}

		_, err := profilesCollection.InsertOne(context.TODO(), user)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true, "message": USERS_ERROR_INSERT_DATABASE_ERROR})
			return
		}

		router.router.ServerConfig.LOGGER.Info("New user profile on validation; " + validateProfileRequest.Wallet)
		c.IndentedJSON(http.StatusOK, user)
		return
	}

	result.Decode(&user)
	c.IndentedJSON(http.StatusOK, user)
}

// SECTION - Router Main methods
// All the methods related to the initialization or configuration
// Normally this methods will be called from another core modules

func (router *UsersRouter) CreateRoutes() error {
	router.router.ParsedGet("/holders", router.getAllHolders)
	router.router.ParsedPost("/possible", router.addPossibleUser)
	router.router.ParsedGet("/possible", router.getPossibleUsers)
	router.router.ParsedPost("/profile", router.createProfileUsingSign)
	router.router.ParsedPost("/profile-validate", router.validateProfile)
	router.router.ParsedGet("/profile/:address", router.getUserProfile)
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
