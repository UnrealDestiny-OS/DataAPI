package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Address string             `bson:"address"`
}

type PossibleUser struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Address    string             `bson:"address"`
	Connection int32              `bson:"connection"`
}

type UserHolder struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Address  string             `bson:"address"`
	Holdings string             `bson:"holdings"`
	Valid    bool               `bson:"valid"`
	Network  string             `bson:"network"`
}

var COLLECTION_HOLDERS = "users_holders"
var COLLECTION_POSSIBLE_USERS = "users_possible"
