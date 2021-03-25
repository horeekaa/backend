package mongodbaccountdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	model "github.com/horeekaa/backend/model"
)

type MemberAccessDataSourceMongo interface {
	FindByID(ID interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error)
	Find(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) ([]*model.MemberAccess, error)
	Create(input *model.CreateMemberAccess, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error)
	Update(ID interface{}, updateData *model.UpdateMemberAccess, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error)
}
