package invoicedomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	databaseinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/interfaces/sources"
	invoicedomainrepositories "github.com/horeekaa/backend/features/invoices/data/repositories"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	databasepaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/interfaces/sources"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
)

type UpdateInvoiceDependency struct{}

func (_ *UpdateInvoiceDependency) Bind() {
	container.Singleton(
		func(
			invoiceDataSource databaseinvoicedatasourceinterfaces.InvoiceDataSource,
			purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
			paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource,
			mouDataSource databasemoudatasourceinterfaces.MouDataSource,
		) invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent {
			updateInvoiceComponent, _ := invoicedomainrepositories.NewUpdateInvoiceTransactionComponent(
				invoiceDataSource,
				purchaseOrderDataSource,
				paymentDataSource,
				mouDataSource,
			)
			return updateInvoiceComponent
		},
	)

	container.Transient(
		func(
			memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
			trxComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent,
			createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) invoicedomainrepositoryinterfaces.UpdateInvoiceRepository {
			updateInvoiceRepo, _ := invoicedomainrepositories.NewUpdateInvoiceRepository(
				memberAccessDataSource,
				trxComponent,
				createNotificationComponent,
				mongoDBTransaction,
			)
			return updateInvoiceRepo
		},
	)
}
