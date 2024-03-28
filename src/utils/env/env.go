package env

import (
	"crypto/ecdsa"
	"log"
	"os"
	"strconv"
	"unrealDestiny/dataAPI/src/utils/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

// NOTE - LoadEnv (*ServerConfig)
// Load all environmental variables using the dotenv library
// All variables should be declares into the Server config struct
// Then the LoadEnv function set all the information in their own variable
// The ServerConfig Env variables should not change over the time
func LoadEnv(serverConfig *config.ServerConfig) bool {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error reading the ENV file.")
		return false
	}

	ENV := os.Getenv("ENV")

	if ENV == "" {
		log.Println("Error getting the ENV variable.")
		return false
	}

	PORT := os.Getenv("PORT")

	if PORT == "" {
		log.Println("Error getting the PORT variable.")
		return false
	}

	MONGO_CLIENT := os.Getenv("MONGO_URI")

	if MONGO_CLIENT == "" {
		log.Println("Error getting the MONGO_URI variable.")
		return false
	}

	MONGO_DATABASE := os.Getenv("MONGO_DATABASE")

	if MONGO_DATABASE == "" {
		log.Println("Error getting the MONGO_DATABASE variable.")
		return false
	}

	MTRG_CLIENT_IP := os.Getenv("MTRG_CLIENT_IP")

	if MTRG_CLIENT_IP == "" {
		log.Println("Error getting the MTRG_CLIENT_IP variable.")
		return false
	}

	MTRG_WS_CLIENT_IP := os.Getenv("MTRG_WS_CLIENT_IP")

	if MTRG_WS_CLIENT_IP == "" {
		log.Println("Error getting the MTRG_WS_CLIENT_IP variable.")
		return false
	}

	USE_PRODUCTION_ADDRESSES := os.Getenv("USE_PRODUCTION_ADDRESSES")

	if USE_PRODUCTION_ADDRESSES == "" {
		log.Println("Error getting the USE_PRODUCTION_ADDRESSES variable.")
		return false
	}

	ACTIVE_CHAIN_ID := os.Getenv("ACTIVE_CHAIN_ID")

	if ACTIVE_CHAIN_ID == "" {
		log.Println("Error getting the ACTIVE_CHAIN_ID variable.")
		return false
	}

	ADMIN_PASS := os.Getenv("ADMIN_PASS")

	if ADMIN_PASS == "" {
		log.Println("Error getting the ADMIN_PASS variable.")
		return false
	}

	ETH_ADDRESS := os.Getenv("ETH_ADDRESS")

	if ETH_ADDRESS == "" {
		log.Println("Error getting the ETH_ADDRESS variable.")
		return false
	}

	var totalPrivateKeys int = 10
	var privateKeysError bool = false
	var loadedPrivateKeys []*ecdsa.PrivateKey

	EXECUTOR_MNEMONIC := os.Getenv("IDLE_GAME_EXECUTOR_NMEMONIC")

	if EXECUTOR_MNEMONIC == "" {
		log.Println("Error getting the EXECUTOR_MNEMONIC variable.")
		return false
	}

	wallet, err := hdwallet.NewFromMnemonic(EXECUTOR_MNEMONIC)

	if err != nil {
		log.Println("Error generating the seed from mnemonic.")
		return false
	}

	for i := 0; i < totalPrivateKeys; i++ {
		path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/" + strconv.Itoa(i))

		account, err := wallet.Derive(path, true)

		if err != nil {
			privateKeysError = true
			break
		}

		log.Println("Loaded new wallet from mnemonic: " + account.Address.String())

		privateKey, err := wallet.PrivateKeyHex(account)

		if err != nil {
			privateKeysError = true
			break
		}

		parsedPrivateKey, err := crypto.HexToECDSA(privateKey)

		if err != nil {
			privateKeysError = true
			break
		}

		loadedPrivateKeys = append(loadedPrivateKeys, parsedPrivateKey)
	}

	if privateKeysError {
		log.Println("Error getting the IDLE_GAME_EXECUTOR_PKEY variables.")
		return false
	}

	serverConfig.ENV = ENV
	serverConfig.PORT = PORT
	serverConfig.MONGO_CLIENT = MONGO_CLIENT
	serverConfig.MONGO_DATABASE = MONGO_DATABASE
	serverConfig.MTRG_CLIENT = MTRG_CLIENT_IP
	serverConfig.MTRG_WS_CLIENT = MTRG_WS_CLIENT_IP
	serverConfig.ADMIN_PASS = ADMIN_PASS
	serverConfig.ETH_ADDRESS = common.HexToAddress(ETH_ADDRESS)

	useProductionAddressParsed, err := strconv.ParseBool(USE_PRODUCTION_ADDRESSES)

	if err != nil {
		log.Println("Error parsing bolean values.")
		return false
	}

	activeChainParsed, err := strconv.Atoi(ACTIVE_CHAIN_ID)

	if err != nil {
		log.Println("Error parsing int values.")
		return false
	}

	serverConfig.USE_PRODUCTION_ADDRESSES = useProductionAddressParsed
	serverConfig.ACTIVE_CHAIN_ID = activeChainParsed
	serverConfig.EXECUTOR_PRIVATE_KEYS = loadedPrivateKeys

	log.Println("Use PRODUCTION_ADDRESSES: " + strconv.FormatBool(serverConfig.USE_PRODUCTION_ADDRESSES))
	log.Println("Use ACTIVE_CHAIN_ID: " + ACTIVE_CHAIN_ID)
	log.Println("Starting Application on " + ENV + " environment")

	return err == nil
}
