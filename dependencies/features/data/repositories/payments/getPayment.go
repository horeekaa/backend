package paymentdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasepaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/interfaces/sources"
	paymentdomainrepositories "github.com/horeekaa/backend/features/payments/data/repositories"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
)

type GetPaymentDependency struct{}

func (_ *GetPaymentDependency) Bind() {
	container.Singleton(
		func(
			paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource,
		) paymentdomainrepositoryinterfaces.GetPaymentRepository {
			getPaymentRepo, _ := paymentdomainrepositories.NewGetPaymentRepository(
				paymentDataSource,
			)
			return getPaymentRepo
		},
	)
}
