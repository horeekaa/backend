package mongodbmemberaccessrefdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	model "github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberAccessRefDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccessRef, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccessRef, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.MemberAccessRef, error)
	Create(input *model.DatabaseCreateMemberAccessRef, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccessRef, error)
	Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdateMemberAccessRef, operationOptions *mongodbcoretypes.OperationOptions) (*model.MemberAccessRef, error)
	GenerateObjectID() primitive.ObjectID
}
