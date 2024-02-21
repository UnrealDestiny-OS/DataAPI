package database

import (
	"unrealDestiny/dataAPI/src/utils/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDatabase(serverConfig *config.ServerConfig) (*mongo.Client, *mongo.Database) {
	client, err := mongo.Connect(serverConfig.CONTEXT, options.Client().ApplyURI(serverConfig.MONGO_CLIENT))

	if err != nil {
		return nil, nil
	}

	serverConfig.LOGGER.Info("Initialized database on " + serverConfig.MONGO_CLIENT + "" + serverConfig.MONGO_DATABASE)

	return client, client.Database(serverConfig.MONGO_DATABASE)
}
