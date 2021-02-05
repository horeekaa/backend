package mongorepointerfaces

import (
	model "github.com/horeekaa/backend/model"
	mongooperations "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations"
)

type PersonRepoMongo interface {
	FindByID(ID interface{}, operationOptions *mongooperations.OperationOptions) (*model.Person, error)
	FindOne(query map[string]interface{}, operationOptions *mongooperations.OperationOptions) (*model.Person, error)
	Find(query map[string]interface{}, operationOptions *mongooperations.OperationOptions) ([]*model.Person, error)
	Create(input *model.CreatePerson, operationOptions *mongooperations.OperationOptions) (*model.Person, error)
	Update(ID interface{}, updateData *model.UpdatePerson, operationOptions *mongooperations.OperationOptions) (*model.Person, error)
}
