package trainers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrainerStatic struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Type       string             `bson:"type"`
	Model      int16              `bson:"model"`
	Experience int32              `bson:"experience"`
	Level      int32              `bson:"level"`
	Health     int32              `bson:"health"`
	Speed      int32              `bson:"speed"`
	Energy     int32              `bson:"energy"`
	Defense    int32              `bson:"defense"`
	Attack     int32              `bson:"attack"`
	Builder    string             `bson:"builder"`
	Name       string             `bson:"name"`
	Network    int16              `bson:"network"`
}

type UserTrainer struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Type       string             `bson:"type"`
	Model      int16              `bson:"model"`
	Experience int32              `bson:"experience"`
	Level      int32              `bson:"level"`
	Health     int32              `bson:"health"`
	Speed      int32              `bson:"speed"`
	Energy     int32              `bson:"energy"`
	Defense    int32              `bson:"defense"`
	Attack     int32              `bson:"attack"`
	Builder    string             `bson:"builder"`
	Name       string             `bson:"name"`
	Network    int16              `bson:"network"`
	Wallet     string             `bson:"wallet"`
	Index      int32              `bson:"index"`
}

var COLLECTION_STATIC_TRAINERS = "static-trainers"
var COLLECTION_USER_TRAINERS = "user-trainers"
