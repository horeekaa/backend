package invoicedomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/interfaces/sources"
	invoicedomainrepositories "github.com/horeekaa/backend/features/invoices/data/repositories"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
)

type GetInvoiceDependency struct{}

func (_ *GetInvoiceDependency) Bind() {
	container.Singleton(
		func(
			invoiceDataSource databaseinvoicedatasourceinterfaces.InvoiceDataSource,
		) invoicedomainrepositoryinterfaces.GetInvoiceRepository {
			getInvoiceRepo, _ := invoicedomainrepositories.NewGetInvoiceRepository(
				invoiceDataSource,
			)
			return getInvoiceRepo
		},
	)
}
