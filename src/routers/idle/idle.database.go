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
}

var COLLECTION_IDLE_LOGS_INJECTIONS = "idle-logs-injection"
