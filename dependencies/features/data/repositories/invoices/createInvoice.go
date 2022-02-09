package invoicedomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databaseinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/interfaces/sources"
	invoicedomainrepositories "github.com/horeekaa/backend/features/invoices/data/repositories"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
)

type CreateInvoiceDependency struct{}

func (_ *CreateInvoiceDependency) Bind() {
	container.Singleton(
		func(
			invoiceDataSource databaseinvoicedatasourceinterfaces.InvoiceDataSource,
			purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
		) invoicedomainrepositoryinterfaces.CreateInvoiceTransactionComponent {
			createInvoiceComponent, _ := invoicedomainrepositories.NewCreateInvoiceTransactionComponent(
				invoiceDataSource,
				purchaseOrderDataSource,
			)
			return createInvoiceComponent
		},
	)

	container.Transient(
		func(
			memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
			trxComponent invoicedomainrepositoryinterfaces.CreateInvoiceTransactionComponent,
			notificationTrx notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) invoicedomainrepositoryinterfaces.CreateInvoiceRepository {
			createInvoiceRepo, _ := invoicedomainrepositories.NewCreateInvoiceRepository(
				memberAccessDataSource,
				trxComponent,
				notificationTrx,
				mongoDBTransaction,
			)
			return createInvoiceRepo
		},
	)
}
