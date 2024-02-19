package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID      primitive.ObjectID `bson:"_id"`
	Address string             `bson:"address"`
}

type UserHolder struct {
	ID       primitive.ObjectID `bson:"_id"`
	Address  string             `bson:"address"`
	Holdings string             `bson:"holdings"`
	Valid    bool               `bson:"valid"`
	Network  int8               `bson:"network"`
}
