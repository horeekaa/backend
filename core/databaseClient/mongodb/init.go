package mongodbcoreclients

import (
	"context"
	"strconv"
	"time"

	mongodbcorewrappers "github.com/horeekaa/backend/core/databaseClient/mongodb/wrappers"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	coreconfigs "github.com/horeekaa/backend/core/commons/configs"
	mongodbcoreclientinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/init"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	horeekaaexceptioncore "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaaexceptioncoreenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
)

// MongoClient holds the database reference to each of the Client collection
type mongoClient struct {
	client       *mongo.Client
	databaseName string
	timeout      time.Duration
}

func newMongoClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, horeekaaexceptioncore.NewExceptionObject(
			horeekaaexceptioncoreenums.ClientInitializationFailed,
			"/newMongoClient",
			err,
		)
	}

	// confirm connection by pinging to primary cluster
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, horeekaaexceptioncore.NewExceptionObject(
			horeekaaexceptioncoreenums.DBConnectionFailed,
			"/newMongoClient",
			err,
		)
	}

	return client, nil
}

func (mongoClient *mongoClient) Connect() (bool, error) {
	timeout, err := strconv.Atoi(coreconfigs.GetEnvVariable(coreconfigs.DbConfigTimeout))
	client, err := newMongoClient(
		coreconfigs.GetEnvVariable(coreconfigs.DbConfigURL),
		timeout,
	)
	if err != nil {
		return false, err
	}

	mongoClient.client = client
	mongoClient.databaseName = coreconfigs.GetEnvVariable(coreconfigs.DbConfigDBName)
	mongoClient.timeout = time.Duration(timeout) * time.Second

	return true, nil
}

func (mongoClient *mongoClient) GetDatabaseName() (string, error) {
	if &mongoClient.databaseName == nil {
		return "", horeekaaexceptioncore.NewExceptionObject(
			horeekaaexceptioncoreenums.ClientInitializationFailed,
			"/newMongoClient",
			nil,
		)
	}
	return mongoClient.databaseName, nil
}

func (mongoClient *mongoClient) GetDatabaseTimeout() (time.Duration, error) {
	if &mongoClient.timeout == nil {
		return time.Duration(0), horeekaaexceptioncore.NewExceptionObject(
			horeekaaexceptioncoreenums.ClientInitializationFailed,
			"/newMongoClient",
			nil,
		)
	}
	return mongoClient.timeout, nil
}

func (mongoClient *mongoClient) GetCollectionRef(collectionName string) (mongodbcorewrapperinterfaces.MongoCollectionRef, error) {
	if mongoClient.client == nil {
		return nil, horeekaaexceptioncore.NewExceptionObject(
			horeekaaexceptioncoreenums.ClientInitializationFailed,
			"/newMongoClient",
			nil,
		)
	}
	colRef := mongoClient.client.Database(mongoClient.databaseName).Collection(collectionName)
	return mongodbcorewrappers.NewMongoCollectionRef(colRef)
}

func (mongoClient *mongoClient) CreateNewSession() (mongodbcorewrapperinterfaces.MongoSession, error) {
	if mongoClient.client == nil {
		return nil, horeekaaexceptioncore.NewExceptionObject(
			horeekaaexceptioncoreenums.ClientInitializationFailed,
			"/newMongoClient",
			nil,
		)
	}
	session, err := mongoClient.client.StartSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (mongoClient *mongoClient) CreateNewSessionContext(
	ctx context.Context,
	sess mongodbcorewrapperinterfaces.MongoSession,
) (mongodbcorewrapperinterfaces.MongoSessionContext, error) {
	return mongo.NewSessionContext(ctx, sess), nil
}

// NewMongoClientRef is getter for the mongodb database reference currently used
func NewMongoClient() (mongodbcoreclientinterfaces.MongoClient, error) {
	return &mongoClient{}, nil
}
