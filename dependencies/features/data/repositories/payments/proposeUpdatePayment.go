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
	paymentdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories/utils"
)

type ProposeUpdatePaymentDependency struct{}

func (_ *ProposeUpdatePaymentDependency) Bind() {
	container.Singleton(
		func(
			paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			paymentDataLoader paymentdomainrepositoryutilityinterfaces.PaymentLoader,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
		) paymentdomainrepositoryinterfaces.ProposeUpdatePaymentTransactionComponent {
			proposeUpdatePaymentComponent, _ := paymentdomainrepositories.NewProposeUpdatePaymentTransactionComponent(
				paymentDataSource,
				loggingDataSource,
				paymentDataLoader,
				mapProcessorUtility,
			)
			return proposeUpdatePaymentComponent
		},
	)

	container.Transient(
		func(
			paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource,
			createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
			proposeUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ProposeUpdateDescriptivePhotoTransactionComponent,
			proposeUpdatePaymentTransactionComponent paymentdomainrepositoryinterfaces.ProposeUpdatePaymentTransactionComponent,
			updateInvoiceTrxComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) paymentdomainrepositoryinterfaces.ProposeUpdatePaymentRepository {
			proposeUpdatePaymentRepo, _ := paymentdomainrepositories.NewProposeUpdatePaymentRepository(
				paymentDataSource,
				createDescriptivePhotoComponent,
				proposeUpdateDescriptivePhotoComponent,
				proposeUpdatePaymentTransactionComponent,
				updateInvoiceTrxComponent,
				mongoDBTransaction,
			)
			return proposeUpdatePaymentRepo
		},
	)
}
