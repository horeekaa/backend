package paymentdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	databasepaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/interfaces/sources"
	paymentdomainrepositories "github.com/horeekaa/backend/features/payments/data/repositories"
	paymentdomainrepositoryinterfaces "github.com/horeekaa/backend/features/payments/domain/repositories"
)

type GetAllPaymentDependency struct{}

func (_ *GetAllPaymentDependency) Bind() {
	container.Singleton(
		func(
			paymentDataSource databasepaymentdatasourceinterfaces.PaymentDataSource,
			mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
		) paymentdomainrepositoryinterfaces.GetAllPaymentRepository {
			getAllPaymentRepo, _ := paymentdomainrepositories.NewGetAllPaymentRepository(
				paymentDataSource,
				mongoQueryBuilder,
			)
			return getAllPaymentRepo
		},
	)
}
