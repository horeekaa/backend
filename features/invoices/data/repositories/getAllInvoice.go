package invoicedomainrepositories

import (
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/interfaces/sources"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	invoicedomainrepositorytypes "github.com/horeekaa/backend/features/invoices/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllInvoiceRepository struct {
	invoiceDataSource databaseinvoicedatasourceinterfaces.InvoiceDataSource
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder
}

func NewGetAllInvoiceRepository(
	invoiceDataSource databaseinvoicedatasourceinterfaces.InvoiceDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (invoicedomainrepositoryinterfaces.GetAllInvoiceRepository, error) {
	return &getAllInvoiceRepository{
		invoiceDataSource,
		mongoQueryBuilder,
	}, nil
}

func (getAllInvoiceRepo *getAllInvoiceRepository) Execute(
	input invoicedomainrepositorytypes.GetAllInvoiceInput,
) ([]*model.Invoice, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllInvoiceRepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	invoices, err := getAllInvoiceRepo.invoiceDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getAllInvoiceRepository",
			err,
		)
	}

	return invoices, nil
}
