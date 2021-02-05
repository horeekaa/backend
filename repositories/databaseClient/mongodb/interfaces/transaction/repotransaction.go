package mongotransactioninterfaces

import (
	mongooperations "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations"
)

type TransactionComponent interface {
	PreTransaction(input interface{}) (interface{}, error)
	TransactionBody(operationOptions *mongooperations.OperationOptions, preOutput interface{}) (interface{}, error)
}

type MongoRepoTransaction interface {
	RunTransaction(input interface{}) (interface{}, error)
}
