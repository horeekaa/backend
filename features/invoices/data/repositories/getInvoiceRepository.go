package invoicedomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/interfaces/sources"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getInvoiceRepository struct {
	invoiceDataSource databaseinvoicedatasourceinterfaces.InvoiceDataSource
	pathIdentity      string
}

func NewGetInvoiceRepository(
	invoiceDataSource databaseinvoicedatasourceinterfaces.InvoiceDataSource,
) (invoicedomainrepositoryinterfaces.GetInvoiceRepository, error) {
	return &getInvoiceRepository{
		invoiceDataSource,
		"GetInvoiceRepository",
	}, nil
}

func (getInvoiceRepo *getInvoiceRepository) Execute(filterFields *model.InvoiceFilterFields) (*model.Invoice, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	invoice, err := getInvoiceRepo.invoiceDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getInvoiceRepo.pathIdentity,
			err,
		)
	}

	return invoice, nil
}
