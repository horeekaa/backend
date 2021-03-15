package mongorepointerfaces

import (
	model "github.com/horeekaa/backend/model"
	mongooperationmodels "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations/models"
)

type MemberAccessRefRepoMongo interface {
	FindByID(ID interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.MemberAccessRef, error)
	FindOne(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.MemberAccessRef, error)
	Find(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) ([]*model.MemberAccessRef, error)
	Create(input *model.CreateMemberAccessRef, operationOptions *mongooperationmodels.OperationOptions) (*model.MemberAccessRef, error)
	Update(ID interface{}, updateData *model.UpdateMemberAccessRef, operationOptions *mongooperationmodels.OperationOptions) (*model.MemberAccessRef, error)
}
