package mongodbaccountdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	model "github.com/horeekaa/backend/model"
)

type AccountDataSourceMongo interface {
	FindByID(ID interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error)
	Find(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) ([]*model.Account, error)
	Create(input *model.CreateAccount, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error)
	Update(ID interface{}, updateData *model.UpdateAccount, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error)
}
