package mongodbcoretransactioninterfaces

import (
	mongodbcoreoperationmodels "github.com/horeekaa/backend/core/databaseClient/mongoDB/operations/models"
)

type TransactionComponent interface {
	PreTransaction(input interface{}) (interface{}, error)
	TransactionBody(operationOptions *mongodbcoreoperationmodels.OperationOptions, preOutput interface{}) (interface{}, error)
}

type MongoRepoTransaction interface {
	RunTransaction(input interface{}) (interface{}, error)
}
