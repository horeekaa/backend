package mongodbcorewrappers

import (
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoInsertOneResult struct {
	*mongo.InsertOneResult
}

func (insertOneRes *mongoInsertOneResult) GetInsertedID() interface{} {
	return insertOneRes.InsertedID
}

func NewMongoInsertOneResult(wrappedInsertOneResult *mongo.InsertOneResult) (mongodbcorewrapperinterfaces.MongoInsertOneResult, error) {
	return &mongoInsertOneResult{
		wrappedInsertOneResult,
	}, nil
}
