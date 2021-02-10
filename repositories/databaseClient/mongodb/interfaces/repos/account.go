package mongorepointerfaces

import (
	model "github.com/horeekaa/backend/model"
	mongooperationmodels "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations/models"
)

type AccountRepoMongo interface {
	FindByID(ID interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.Account, error)
	FindOne(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.Account, error)
	Find(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) ([]*model.Account, error)
	Create(input *model.CreateAccount, operationOptions *mongooperationmodels.OperationOptions) (*model.Account, error)
	Update(ID interface{}, updateData *model.UpdateAccount, operationOptions *mongooperationmodels.OperationOptions) (*model.Account, error)
}
