package mongotransaction

import (
	mongooperations "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations"
)

type MongoTransactionComponent interface {
	preTransaction(input interface{}) (interface{}, error)
	transactionBody(session *mongooperations.OperationOptions, preOutput interface{}) (interface{}, error)
}
