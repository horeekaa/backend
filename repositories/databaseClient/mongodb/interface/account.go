package mongorepointerface

import (
	model "github.com/horeekaa/backend/model"
	mongooperations "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations"
)

type AccountRepoMongo interface {
	FindByID(ID interface{}, operationOptions *mongooperations.OperationOptions) (*model.Account, error)
	FindOne(query mongooperations.OperationQueryType, operationOptions *mongooperations.OperationOptions) (*model.Account, error)
	Find(query mongooperations.OperationQueryType, operationOptions *mongooperations.OperationOptions) ([]*model.Account, error)
	Create(input *model.CreateAccount, operationOptions *mongooperations.OperationOptions) (*model.Account, error)
	Update(ID interface{}, updateData *model.UpdateAccount, operationOptions *mongooperations.OperationOptions) (*model.Account, error)
}
