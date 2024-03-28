package runtime_router

import (
	"fmt"
	"runtime"
	"time"
	"unrealDestiny/dataAPI/src/utils/config"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type RuntimeRouter struct {
	router config.Router
}

// SECTION - REST API
// Rest API methods

// NOTE - AddPossibleUser
// GET Request, No Body, No params
func (router *RuntimeRouter) getDailyLogs(c *gin.Context) {

}

// SECTION - Router Main methods
// All the methods related to the initialization or configuration
// Normally this methods will be called from another core modules

func (router *RuntimeRouter) CreateRoutes() error {
	router.router.ParsedGet("/dailyLogs", router.getDailyLogs)
	return nil
}

func (router *RuntimeRouter) collectRuntimeData() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	runtimeCollection := router.router.MainDatabase.Collection(COLLECTION_RUNTIME_LOGS)

	var log RuntimeLog

	log.CPUNumberOfCPUs = runtime.NumCPU()
	log.CPUNumberOfGoroutines = runtime.NumGoroutine()
	log.MemAlloc = m.Alloc / 1024 / 1024
	log.MemTotalAlloc = m.TotalAlloc / 1024 / 1024
	log.MemSys = m.Sys / 1024 / 1024
	log.MemNumGC = m.NumGC
	log.Date = time.Now().UnixMilli()

	_, err := runtimeCollection.InsertOne(context.Background(), log)

	if err != nil {
		router.router.ServerConfig.LOGGER.Error(RUNTIME_ADD_LOG_ERROR)
		return
	}

	fmt.Printf("Alloc = %v MiB", log.MemAlloc)
	fmt.Printf("\tTotalAlloc = %v MiB", log.MemTotalAlloc)
	fmt.Printf("\tSys = %v MiB", log.MemSys)
	fmt.Printf("\tNumGC = %v\n", log.MemNumGC)
	fmt.Printf("Number of CPUs: %d\n", log.CPUNumberOfCPUs)
	fmt.Printf("Number of goroutines: %d\n", log.CPUNumberOfGoroutines)
}

func (router *RuntimeRouter) StartRuntimeLoop() {
	go func() {
		for {
			time.Sleep(time.Second * 10)
			router.collectRuntimeData()
		}
	}()
}

func (router *RuntimeRouter) Init(serverConfig *config.ServerConfig, mainRouter *gin.Engine, database *mongo.Database) {
	router.router.ServerConfig = serverConfig
	router.router.MainRouter = mainRouter
	router.router.MainDatabase = database
	router.router.Name = "runtime"
	router.router.Path = "/runtime"
	router.router.ServerConfig.LOGGER.Info("Starting Runtime router on " + router.router.Path)
	router.StartRuntimeLoop()
}

func CreateRouter(serverConfig *config.ServerConfig, router *gin.Engine, database *mongo.Database) *RuntimeRouter {
	users := new(RuntimeRouter)
	users.Init(serverConfig, router, database)
	return users
}
