package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Address string             `bson:"address"`
}

type UserHolder struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Address  string             `bson:"address"`
	Holdings string             `bson:"holdings"`
	Valid    bool               `bson:"valid"`
	Network  string             `bson:"network"`
}

type UserHolderData struct {
	Address  string `json:"address,omitempty"`
	Name     string `json:"name,omitempty"`
	Holdings string `json:"amount,omitempty"`
	Valid    string `json:"valid,omitempty"`
	Network  string `json:"network,omitempty"`
}
