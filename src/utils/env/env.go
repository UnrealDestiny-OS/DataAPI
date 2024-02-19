package env

import (
	"log"
	"os"
	"unrealDestiny/dataAPI/src/utils/config"

	"github.com/joho/godotenv"
)

// NOTE - LoadEnv (*ServerConfig)
// Load all environmental variables using the dotenv library
// All variables should be declares into the Server config struct
// Then the LoadEnv function set all the information in their own variable
// The ServerConfig Env variables should not change over the time
func LoadEnv(serverConfig *config.ServerConfig) bool {
	err := godotenv.Load(".env")

	ENV := os.Getenv("ENV")

	if ENV == "" {
		return false
	}

	PORT := os.Getenv("PORT")

	if PORT == "" {
		return false
	}

	MONGO_CLIENT := os.Getenv("MONGO_URI")

	if MONGO_CLIENT == "" {
		return false
	}

	MTRG_CLIENT_IP := os.Getenv("MTRG_CLIENT_IP")

	if MTRG_CLIENT_IP == "" {
		return false
	}

	serverConfig.ENV = ENV
	serverConfig.PORT = PORT
	serverConfig.MONGO_CLIENT = MONGO_CLIENT
	serverConfig.MTRG_CLIENT = MTRG_CLIENT_IP

	log.Println("Starting Application on " + ENV + " environment")

	return err == nil
}
