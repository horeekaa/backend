package mongodbcoretransactioninterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionComponent interface {
	PreTransaction(input interface{}) (interface{}, error)
	TransactionBody(operationOptions *mongodbcoretypes.OperationOptions, preOutput interface{}) (interface{}, error)
}

type MongoRepoTransaction interface {
	SetTransaction(component TransactionComponent, transactionTitle string) bool
	RunTransaction(input interface{}) (interface{}, error)
	TransactionFn(
		preTransactOutput interface{},
	) func(sessCtx mongo.SessionContext) (interface{}, error)
}
