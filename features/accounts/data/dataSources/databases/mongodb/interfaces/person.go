package mongodbaccountdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	model "github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PersonDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.Person, error)
	Create(input *model.DatabaseCreatePerson, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error)
	Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdatePerson, operationOptions *mongodbcoretypes.OperationOptions) (*model.Person, error)
}
