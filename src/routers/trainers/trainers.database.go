package trainers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrainerStatic struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}
