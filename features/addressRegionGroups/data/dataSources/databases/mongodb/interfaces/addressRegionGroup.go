package mongodbaddressregiongroupdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	model "github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddressRegionGroupDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.AddressRegionGroup, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.AddressRegionGroup, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.AddressRegionGroup, error)
	Create(input *model.DatabaseCreateAddressRegionGroup, operationOptions *mongodbcoretypes.OperationOptions) (*model.AddressRegionGroup, error)
	Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdateAddressRegionGroup, operationOptions *mongodbcoretypes.OperationOptions) (*model.AddressRegionGroup, error)
	GenerateObjectID() primitive.ObjectID
}
