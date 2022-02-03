package invoicedomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	databaseinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/interfaces/sources"
	invoicedomainrepositories "github.com/horeekaa/backend/features/invoices/data/repositories"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
)

type GetAllInvoiceDependency struct{}

func (_ *GetAllInvoiceDependency) Bind() {
	container.Singleton(
		func(
			invoiceDataSource databaseinvoicedatasourceinterfaces.InvoiceDataSource,
			mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
		) invoicedomainrepositoryinterfaces.GetAllInvoiceRepository {
			getAllInvoiceRepo, _ := invoicedomainrepositories.NewGetAllInvoiceRepository(
				invoiceDataSource,
				mongoQueryBuilder,
			)
			return getAllInvoiceRepo
		},
	)
}
