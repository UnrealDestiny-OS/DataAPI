package database

import (
	"unrealDestiny/dataAPI/src/utils/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDatabase(serverConfig *config.ServerConfig) *mongo.Client {
	client, err := mongo.Connect(serverConfig.CONTEXT, options.Client().ApplyURI(serverConfig.MONGO_CLIENT))

	if err != nil {
		return nil
	}

	defer client.Disconnect(serverConfig.CONTEXT)

	serverConfig.LOGGER.Info("Initialized database on " + serverConfig.MONGO_CLIENT)

	return client
}
