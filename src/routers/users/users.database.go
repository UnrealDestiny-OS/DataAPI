package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStaticPossible struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Address    string             `bson:"address"`
	Connection int32              `bson:"connection"`
}

type UserStaticHolder struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Address  string             `bson:"address"`
	Holdings string             `bson:"holdings"`
	Valid    bool               `bson:"valid"`
	Network  string             `bson:"network"`
}

var COLLECTION_HOLDERS = "static-users-holders"
var COLLECTION_POSSIBLE_USERS = "static-users-possible"
