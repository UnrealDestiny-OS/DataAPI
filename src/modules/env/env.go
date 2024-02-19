package env

import (
	"log"
	"os"
	"unrealDestiny/contractsReader/src/modules/config"

	"github.com/joho/godotenv"
)

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

	MTRG_CLIENT_IP := os.Getenv("MTRG_CLIENT_IP")

	if MTRG_CLIENT_IP == "" {
		return false
	}

	serverConfig.ENV = ENV
	serverConfig.PORT = PORT
	serverConfig.MTRG_CLIENT = MTRG_CLIENT_IP

	log.Println("Starting Application on " + ENV + " environment")

	return err == nil
}
