package mongodbmemberaccessdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	model "github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberAccessDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.MemberAccess, error)
	Create(input *model.CreateMemberAccess, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error)
	Update(ID primitive.ObjectID, updateData *model.UpdateMemberAccess, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccess, error)
}
