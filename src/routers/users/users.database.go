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

type UserProfile struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Wallet     string             `bson:"wallet"`
	Name       string             `bson:"name"`
	UDT        int64              `bson:"UDT"`
	FUDT       int64              `bson:"FUDT"`
	Level      int64              `bson:"level"`
	Experience int64              `bson:"experience"`
	Chain      int64              `bson:"chain"`
}

// SECTION -API Requests Types

type APICreateUserProfile struct {
	Username      string `json:"username"`
	Wallet        string `json:"wallet"`
	Sign          string `json:"sign"`
	CreationChain int64  `json:"chain"`
}

var COLLECTION_HOLDERS = "static-users-holders"
var COLLECTION_POSSIBLE_USERS = "static-users-possible"
var COLLECTION_USER_PROFILES = "user-profiles"
