package paymentdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasepaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/interfaces/sources"
	paymentdomainrepositories "github.com/horeekaa/backend/features/payments/data/repositories"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
)

type ApproveUpdatePaymentDependency struct{}

func (_ *ApproveUpdatePaymentDependency) Bind() {
	container.Singleton(
		func(
			paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) paymentdomainrepositoryinterfaces.ApproveUpdatePaymentTransactionComponent {
			approveUpdatePaymentComponent, _ := paymentdomainrepositories.NewApproveUpdatePaymentTransactionComponent(
				paymentDataSource,
				loggingDataSource,
				mapProcessorUtility,
			)
			return approveUpdatePaymentComponent
		},
	)

	container.Transient(
		func(
			paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource,
			approveUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent,
			trxComponent paymentdomainrepositoryinterfaces.ApproveUpdatePaymentTransactionComponent,
			updateInvoiceTrxComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) paymentdomainrepositoryinterfaces.ApproveUpdatePaymentRepository {
			approveUpdatePaymentRepo, _ := paymentdomainrepositories.NewApproveUpdatePaymentRepository(
				paymentDataSource,
				approveUpdateDescriptivePhotoComponent,
				trxComponent,
				updateInvoiceTrxComponent,
				mongoDBTransaction,
			)
			return approveUpdatePaymentRepo
		},
	)
}
