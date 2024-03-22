package idle

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IdleInjectionLog struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Wallet  string             `bson:"wallet"`
	Amount  string             `bson:"amount"`
	Date    string             `bson:"date"`
	Network int                `bson:"network"`
	Type    string             `bson:"type"`
}

var COLLECTION_IDLE_LOGS_INJECTIONS = "idle-logs-injection"

//SECTION - Api requests typing

type APIExecuteTrainerJoin struct {
	Wallet     string `json:"wallet"`
	Trainer    int    `json:"trainer"`
	WalletAuth string `json:"walletAuth"`
}

type APIExecuteCollectTransactionPoints struct {
	Wallet     string `json:"wallet"`
	Trainer    int    `json:"trainer"`
	WalletAuth string `json:"walletAuth"`
}
