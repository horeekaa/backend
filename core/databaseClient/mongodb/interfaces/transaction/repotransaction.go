package mongodbcoretransactioninterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
)

type TransactionComponent interface {
	PreTransaction(input interface{}) (interface{}, error)
	TransactionBody(operationOptions *mongodbcoretypes.OperationOptions, preOutput interface{}) (interface{}, error)
}

type MongoRepoTransaction interface {
	SetTransaction(component TransactionComponent, transactionTitle string) bool
	RunTransaction(input interface{}) (interface{}, error)
}
