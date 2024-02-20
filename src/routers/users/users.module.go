package users

import (
	"context"
	"fmt"
	"net/http"
	"unrealDestiny/dataAPI/src/utils/config"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

// NOTE - ParsedPost
func (router *UsersRouter) ParsedPost(path string, callback func(*gin.Context)) {
	router.router.MainRouter.POST(router.router.Path+path, callback)
}

// SECTION - REST API
// Rest API methods

func (router *UsersRouter) GetAllUsers(context *gin.Context) {

}

// NOTE - AddPossibleUser
// GET Request, No Body, No params
func (router *UsersRouter) GetPossibleUsers(c *gin.Context) {
	possibleUsersCollection := router.router.MainDatabase.Collection(COLLECTION_POSSIBLE_USERS)

	var users []PossibleUser

	cursor, err := possibleUsersCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
		return
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, users)
}

// NOTE - AddPossibleUser
// POST Request, Body *PossibleUser
// Insert new possible user when the user reach the web page and connect the wallet to the site
func (router *UsersRouter) AddPossibleUser(c *gin.Context) {
	var user PossibleUser
	var searchedUser PossibleUser

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
		fmt.Print(err)
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

// NOTE - GetAllHolders
// GET Request, No Body, No params
// Return all users data on the holders database collection
func (router *UsersRouter) GetAllHolders(c *gin.Context) {
	holdersCollection := router.router.MainDatabase.Collection(COLLECTION_HOLDERS)

	var users []UserHolder

	cursor, err := holdersCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
		return
	}

	if err = cursor.All(context.TODO(), &users); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, users)
}

// SECTION - Router Main methods
// All the methods related to the initialization or configuration
// Normally this methods will be called from another core modules

func (router *UsersRouter) CreateRoutes() error {
	router.ParsedGet("/all", router.GetAllUsers)
	router.ParsedGet("/holders", router.GetAllHolders)
	router.ParsedPost("/possible", router.AddPossibleUser)
	router.ParsedGet("/possible", router.GetPossibleUsers)
	// router.ParsedPost("/holders", router.UploadAllHolders)
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
