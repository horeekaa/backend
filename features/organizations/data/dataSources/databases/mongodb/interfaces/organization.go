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
	Create(input *model.DatabaseCreateOrganization, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error)
	Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdateOrganization, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error)
	GenerateObjectID() primitive.ObjectID
}
