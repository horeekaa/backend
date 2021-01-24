package mongotransaction

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTransactionComponent interface {
	preTransaction(input interface{}) (interface{}, error)
	transactionBody(session *mongo.SessionContext, preOutput interface{}) (interface{}, error)
}
