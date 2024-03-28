package config

import (
	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type ServerConfig struct {
	ENV                      string
	PORT                     string
	MONGO_CLIENT             string
	MONGO_DATABASE           string
	MTRG_CLIENT              string
	MTRG_WS_CLIENT           string
	LOGGER                   *zap.Logger
	CONTEXT                  context.Context
	USE_PRODUCTION_ADDRESSES bool
	ACTIVE_CHAIN_ID          int
	ADMIN_PASS               string
	ETH_ADDRESS              common.Address
	EXECUTOR_PRIVATE_KEYS    []*ecdsa.PrivateKey
}

type Router struct {
	Name         string
	Path         string
	ServerConfig *ServerConfig
	MainRouter   *gin.Engine
	MainDatabase *mongo.Database
	ETHCLient    *ethclient.Client
}

// NOTE - ParsedGet
// Modify the initial Get function to add the router initial Path, in this case /<InitialRouterPath>/<NewRoutePath>
func (router *Router) ParsedGet(path string, callback func(*gin.Context)) {
	router.MainRouter.GET(router.Path+path, callback)
}

// NOTE - ParsedPost
func (router *Router) ParsedPost(path string, callback func(*gin.Context)) {
	router.MainRouter.POST(router.Path+path, callback)
}
