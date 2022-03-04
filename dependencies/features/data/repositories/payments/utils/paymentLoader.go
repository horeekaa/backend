package paymentdomainrepositoryutilitydependencies

import (
	"github.com/golobby/container/v2"
	databaseinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/interfaces/sources"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	paymentdomainrepositoryutilities "github.com/horeekaa/backend/features/payments/data/repositories/utils"
	paymentdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories/utils"
)

type PaymentLoaderDependency struct{}

func (_ *PaymentLoaderDependency) Bind() {
	container.Singleton(
		func(
			invoiceDataSource databaseinvoicedatasourceinterfaces.InvoiceDataSource,
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
		) paymentdomainrepositoryutilityinterfaces.PaymentLoader {
			paymentLoader, _ := paymentdomainrepositoryutilities.NewPaymentLoader(
				invoiceDataSource,
				organizationDataSource,
			)
			return paymentLoader
		},
	)
}
