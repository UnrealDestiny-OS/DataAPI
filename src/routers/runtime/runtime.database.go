package runtime_router

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RuntimeLog struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty"`
	Date                  int64              `bson:"Date"`
	MemAlloc              uint64             `bson:"MemAlloc"`
	MemTotalAlloc         uint64             `bson:"MemTotalAlloc"`
	MemSys                uint64             `bson:"MemSys"`
	MemNumGC              uint32             `bson:"MemNumGC"`
	CPUNumberOfCPUs       int                `bson:"CPUNumberOfCPUs"`
	CPUNumberOfGoroutines int                `bson:"CPUNumberOfGoroutines"`
}

// SECTION - API Requests Types

type APIGetDailyLogsRequest struct {
}

var COLLECTION_RUNTIME_LOGS = "runtime-logs"
