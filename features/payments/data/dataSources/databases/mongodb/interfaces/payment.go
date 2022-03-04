package mongodbpaymentdatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Payment, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Payment, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.Payment, error)
	Create(input *model.DatabaseCreatePayment, operationOptions *mongodbcoretypes.OperationOptions) (*model.Payment, error)
	Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdatePayment, operationOptions *mongodbcoretypes.OperationOptions) (*model.Payment, error)
	GenerateObjectID() primitive.ObjectID
}
