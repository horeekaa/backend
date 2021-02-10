package mongotransactioninterfaces

import (
	mongooperationmodels "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations/models"
)

type TransactionComponent interface {
	PreTransaction(input interface{}) (interface{}, error)
	TransactionBody(operationOptions *mongooperationmodels.OperationOptions, preOutput interface{}) (interface{}, error)
}

type MongoRepoTransaction interface {
	RunTransaction(input interface{}) (interface{}, error)
}
