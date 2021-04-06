package mongodbcoreclientinterfaces

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoClient interface {
	Connect() (bool, error)
	GetClient() (*mongo.Client, error)
	GetDatabaseName() (string, error)
	GetDatabaseTimeout() (time.Duration, error)
}
