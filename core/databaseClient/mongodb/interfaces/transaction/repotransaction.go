package mongodbcoretransactioninterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
)

type TransactionComponent interface {
	PreTransaction(input interface{}) (interface{}, error)
	TransactionBody(operationOptions *mongodbcoretypes.OperationOptions, preOutput interface{}) (interface{}, error)
}

type MongoRepoTransaction interface {
	SetTransaction(component interface{}, transactionTitle string) bool
	RunTransaction(input interface{}) (interface{}, error)
}
