package mongorepointerfaces

import (
	model "github.com/horeekaa/backend/model"
	mongooperationmodels "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations/models"
)

type PersonRepoMongo interface {
	FindByID(ID interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.Person, error)
	FindOne(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.Person, error)
	Find(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) ([]*model.Person, error)
	Create(input *model.CreatePerson, operationOptions *mongooperationmodels.OperationOptions) (*model.Person, error)
	Update(ID interface{}, updateData *model.UpdatePerson, operationOptions *mongooperationmodels.OperationOptions) (*model.Person, error)
}
