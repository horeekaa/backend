package invoicedomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databaseinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/interfaces/sources"
	invoicedomainrepositories "github.com/horeekaa/backend/features/invoices/data/repositories"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
)

type UpdateInvoiceDependency struct{}

func (_ *UpdateInvoiceDependency) Bind() {
	container.Singleton(
		func(
			invoiceDataSource databaseinvoicedatasourceinterfaces.InvoiceDataSource,
			purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
		) invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent {
			updateInvoiceComponent, _ := invoicedomainrepositories.NewUpdateInvoiceTransactionComponent(
				invoiceDataSource,
				purchaseOrderDataSource,
			)
			return updateInvoiceComponent
		},
	)

	container.Transient(
		func(
			trxComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) invoicedomainrepositoryinterfaces.UpdateInvoiceRepository {
			updateInvoiceRepo, _ := invoicedomainrepositories.NewUpdateInvoiceRepository(
				trxComponent,
				mongoDBTransaction,
			)
			return updateInvoiceRepo
		},
	)
}
