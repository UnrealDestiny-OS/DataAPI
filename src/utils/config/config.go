package config

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type ServerConfig struct {
	ENV          string
	PORT         string
	MONGO_CLIENT string
	MTRG_CLIENT  string
	LOGGER       *zap.Logger
	CONTEXT      context.Context
}

type Router struct {
	Name         string
	Path         string
	ServerConfig *ServerConfig
	MainRouter   *gin.Engine
	MainDatabase *mongo.Database
}
