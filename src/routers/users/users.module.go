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
// POST Request, Body
//
//	{
//		"address": "0x",
//	 	"connection": 100000
//	}
//
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

// NOTE - Upload all holders
// One time method, it will include all validated holders on the database
// Holder format
//
//	{
//		"address": "0xA5a01eE4809aA9fC66B9e8dAF20ea7A6466aa8D2",
//		"name": "Hacker",
//		"valid": "INVALID",
//		"amount": "57487079.582",
//		"network": "Polygon"
//	}
//
// It should conver the INVALID string format to boolean format
// func (router *UsersRouter) UploadAllHolders(ginContext *gin.Context) {
// 	raw, err := ginContext.GetRawData()
//
// 	if err != nil {
// 		ginContext.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
// 		return
// 	}
//
// 	var buf bytes.Buffer
//
// 	if err := json.Compact(&buf, raw); err != nil {
// 		ginContext.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
// 		return
// 	}
//
// 	in := []byte(buf.Bytes())
// 	holders := []UserHolderData{}
//
// 	err = json.Unmarshal(in, &holders)
//
// 	if err != nil {
// 		ginContext.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
// 		return
// 	}
//
// 	if len(holders) > 0 {
// 		holdersCollection := router.router.MainDatabase.Collection(HOLDERS_COLLECTION)
// 		parsedHolders := []interface{}{}
//
// 		for i := 0; i < len(holders); i++ {
// 			parsedHolders = append(parsedHolders, UserHolder{
// 				Address:  holders[i].Address,
// 				Holdings: holders[i].Holdings,
// 				Valid:    holders[i].Valid == "VALID",
// 				Network:  holders[i].Address,
// 			})
// 		}
//
// 		_, err := holdersCollection.InsertMany(context.TODO(), parsedHolders)
//
// 		if err != nil {
// 			fmt.Println(err)
// 			ginContext.IndentedJSON(http.StatusBadRequest, gin.H{"error": true})
// 		}
// 	}
//
// 	ginContext.IndentedJSON(http.StatusOK, nil)
//
// }

// SECTION - Router Main methods
// All the methods related to the initialization or configuration
// Normally this methods will be called from another core modules

func (router *UsersRouter) CreateRoutes() error {
	router.ParsedGet("/all", router.GetAllUsers)
	router.ParsedGet("/holders/all", router.GetAllHolders)
	router.ParsedPost("/possible", router.AddPossibleUser)
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
