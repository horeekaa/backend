package mongodborganizationdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongoDB/types"
	"github.com/horeekaa/backend/model"
)

type OrganizationDataSourceMongo interface {
	FindByID(ID interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error)
	Find(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) ([]*model.Organization, error)
	Create(input *model.CreateOrganization, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error)
	Update(ID interface{}, updateData *model.UpdateOrganization, operationOptions *mongodbcoretypes.OperationOptions) (*model.Organization, error)
}
