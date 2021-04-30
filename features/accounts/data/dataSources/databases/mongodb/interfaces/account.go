package mongodbaccountdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	model "github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.Account, error)
	Create(input *model.CreateAccount, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error)
	Update(ID primitive.ObjectID, updateData *model.UpdateAccount, operationOptions *mongodbcoretypes.OperationOptions) (*model.Account, error)
}
