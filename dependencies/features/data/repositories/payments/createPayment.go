package paymentdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasepaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/interfaces/sources"
	paymentdomainrepositories "github.com/horeekaa/backend/features/payments/data/repositories"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
	paymentdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories/utils"
)

type CreatePaymentDependency struct{}

func (_ *CreatePaymentDependency) Bind() {
	container.Singleton(
		func(
			paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			paymentDataLoader paymentdomainrepositoryutilityinterfaces.PaymentLoader,
			updateInvoiceTrxComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent,
		) paymentdomainrepositoryinterfaces.CreatePaymentTransactionComponent {
			createPaymentComponent, _ := paymentdomainrepositories.NewCreatePaymentTransactionComponent(
				paymentDataSource,
				loggingDataSource,
				paymentDataLoader,
				updateInvoiceTrxComponent,
			)
			return createPaymentComponent
		},
	)

	container.Transient(
		func(
			trxComponent paymentdomainrepositoryinterfaces.CreatePaymentTransactionComponent,
			createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) paymentdomainrepositoryinterfaces.CreatePaymentRepository {
			createPaymentRepo, _ := paymentdomainrepositories.NewCreatePaymentRepository(
				trxComponent,
				createDescriptivePhotoComponent,
				mongoDBTransaction,
			)
			return createPaymentRepo
		},
	)
}
