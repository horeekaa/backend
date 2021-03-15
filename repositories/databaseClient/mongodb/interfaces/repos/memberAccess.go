package mongorepointerfaces

import (
	model "github.com/horeekaa/backend/model"
	mongooperationmodels "github.com/horeekaa/backend/repositories/databaseClient/mongoDB/operations/models"
)

type MemberAccessRepoMongo interface {
	FindByID(ID interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.MemberAccess, error)
	FindOne(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) (*model.MemberAccess, error)
	Find(query map[string]interface{}, operationOptions *mongooperationmodels.OperationOptions) ([]*model.MemberAccess, error)
	Create(input *model.CreateMemberAccess, operationOptions *mongooperationmodels.OperationOptions) (*model.MemberAccess, error)
	Update(ID interface{}, updateData *model.UpdateMemberAccess, operationOptions *mongooperationmodels.OperationOptions) (*model.MemberAccess, error)
}
