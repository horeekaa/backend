package mongodborganizationdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrganizationDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.Organization, error)
	Create(input *model.InternalCreateOrganization, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error)
	Update(ID primitive.ObjectID, updateData *model.InternalUpdateOrganization, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error)
}
