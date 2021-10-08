package mongodbAddressdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	model "github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddressDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Address, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Address, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.Address, error)
	Create(input *model.DatabaseCreateAddress, operationOptions *mongodbcoretypes.OperationOptions) (*model.Address, error)
	Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdateAddress, operationOptions *mongodbcoretypes.OperationOptions) (*model.Address, error)
}
