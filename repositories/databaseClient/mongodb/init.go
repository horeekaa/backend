package mongodbclients

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	horeekaaexception "github.com/horeekaa/backend/_errors/repoExceptions"
	horeekaaexceptionenums "github.com/horeekaa/backend/_errors/repoExceptions/_enums"
)

var (
	DatabaseClient *MongoRepository
)

// MongoRepository holds the database reference to each of the repository collection
type MongoRepository struct {
	Client       *mongo.Client
	DatabaseName string
	Timeout      time.Duration
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.DBConnectionFailed,
			"/newMongoClient",
			err,
		)
	}

	// confirm connection by pinging to primary cluster
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.DBConnectionFailed,
			"/newMongoClient",
			err,
		)
	}

	return client, nil
}

// NewMongoClientRef is getter for the mongodb database reference currently used
func NewMongoClientRef(mongoURL string, databaseName string, mongoTimeout int) (*MongoRepository, error) {
	client, err := newMongoClient(mongoURL, mongoTimeout)
	if err != nil {
		return nil, err
	}

	return &MongoRepository{
		Client:       client,
		DatabaseName: databaseName,
		Timeout:      time.Duration(mongoTimeout) * time.Second,
	}, nil
}
