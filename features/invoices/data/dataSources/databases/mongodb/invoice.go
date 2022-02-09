package mongodbinvoicedatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type invoiceDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewInvoiceDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbinvoicedatasourceinterfaces.InvoiceDataSourceMongo, error) {
	basicOperation.SetCollection("invoices")
	return &invoiceDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
}

func (invoiceDataSourceMongo *invoiceDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (invoiceDataSourceMongo *invoiceDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Invoice, error) {
	var output model.Invoice
	_, err := invoiceDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (invoiceDataSourceMongo *invoiceDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Invoice, error) {
	var output model.Invoice
	_, err := invoiceDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (invoiceDataSourceMongo *invoiceDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.Invoice, error) {
	var invoices = []*model.Invoice{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var invoice model.Invoice
		if err := cursor.Decode(&invoice); err != nil {
			return err
		}
		invoices = append(invoices, &invoice)
		return nil
	}
	_, err := invoiceDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return invoices, err
}

func (invoiceDataSourceMongo *invoiceDataSourceMongo) Create(input *model.DatabaseCreateInvoice, operationOptions *mongodbcoretypes.OperationOptions) (*model.Invoice, error) {
	_, err := invoiceDataSourceMongo.setDefaultValuesWhenCreate(
		input,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.Invoice
	_, err = invoiceDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (invoiceDataSourceMongo *invoiceDataSourceMongo) Update(updateCriteria map[string]interface{}, updateData *model.DatabaseUpdateInvoice, operationOptions *mongodbcoretypes.OperationOptions) (*model.Invoice, error) {
	_, err := invoiceDataSourceMongo.setDefaultValuesWhenUpdate(
		updateCriteria,
		updateData,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.Invoice
	_, err = invoiceDataSourceMongo.basicOperation.Update(
		updateCriteria,
		map[string]interface{}{
			"$set": updateData,
		},
		&output,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (invoiceDataSourceMongo *invoiceDataSourceMongo) setDefaultValuesWhenUpdate(
	inputCriteria map[string]interface{},
	input *model.DatabaseUpdateInvoice,
	operationOptions *mongodbcoretypes.OperationOptions,
) (bool, error) {
	currentTime := time.Now()
	existingObject, err := invoiceDataSourceMongo.FindOne(inputCriteria, operationOptions)
	if err != nil {
		return false, err
	}
	if existingObject == nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.QueryObjectFailed,
			"/invoiceDataSource/update",
			nil,
		)
	}
	input.UpdatedAt = &currentTime

	return true, nil
}

func (invoiceDataSourceMongo *invoiceDataSourceMongo) setDefaultValuesWhenCreate(
	input *model.DatabaseCreateInvoice,
) (bool, error) {
	currentTime := time.Now()
	defaultStatus := model.InvoiceStatusAvailable

	if input.Status == nil {
		input.Status = &defaultStatus
	}
	input.CreatedAt = &currentTime
	input.UpdatedAt = &currentTime

	return true, nil
}
