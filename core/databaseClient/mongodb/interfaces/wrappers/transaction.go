package mongodbcorewrapperinterfaces

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoSession interface {
	mongo.Session
}

type MongoSessionContext interface {
	mongo.SessionContext
}
