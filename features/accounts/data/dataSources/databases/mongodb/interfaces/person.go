package mongodbaccountdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	model "github.com/horeekaa/backend/model"
)

type PersonRepoMongo interface {
	FindByID(ID interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error)
	Find(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) ([]*model.Person, error)
	Create(input *model.CreatePerson, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error)
	Update(ID interface{}, updateData *model.UpdatePerson, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error)
}
