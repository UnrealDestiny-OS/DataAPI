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

type TransactionExecutionLog struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	GenerationID string             `bson:"generationID"`
	Hash         string             `bson:"hash"`
	Error        string             `bson:"error"`
	Network      int                `bson:"network"`
}

var COLLECTION_IDLE_LOGS_INJECTIONS = "idle-logs-injection"
var COLLECTION_IDLE_EXECUTION_LOG = "idle-logs-execution"

//SECTION - Api requests typing

type APIExecuteTrainerJoin struct {
	Wallet     string `json:"wallet"`
	Trainer    int    `json:"trainer"`
	WalletAuth string `json:"walletAuth"`
	Chain      int    `json:"chain"`
}

type APIExecuteCollectTransactionPoints struct {
	Wallet     string `json:"wallet"`
	Trainer    int    `json:"trainer"`
	WalletAuth string `json:"walletAuth"`
	Chain      int    `json:"chain"`
}

type APIExecuteImprovementBuy struct {
	Wallet      string `json:"wallet"`
	Trainer     int    `json:"trainer"`
	Improvement int    `json:"improvement"`
	WalletAuth  string `json:"walletAuth"`
	Chain       int    `json:"chain"`
}

type APIAdminRequest struct {
	AdminPassword    string `json:"adminPassword"`
	AccountInbalance string `json:"accountInbalance"`
}
