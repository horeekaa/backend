package mongodbinvoicedatasourceinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvoiceDataSourceMongo interface {
	FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Invoice, error)
	FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Invoice, error)
	Find(
		query map[string]interface{},
		paginationOpts *mongodbcoretypes.PaginationOptions,
		operationOptions *mongodbcoretypes.OperationOptions,
	) ([]*model.Invoice, error)
	Create(input *model.DatabaseCreateInvoice, operationOptions *mongodbcoretypes.OperationOptions) (*model.Invoice, error)
	Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdateInvoice, operationOptions *mongodbcoretypes.OperationOptions) (*model.Invoice, error)
	GenerateObjectID() primitive.ObjectID
}
