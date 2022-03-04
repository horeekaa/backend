package mongodbpaymentdatasourcedependencies

import (
	container "github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databasepaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/interfaces/sources"
	mongodbpaymentdatasources "github.com/horeekaa/backend/features/payments/data/dataSources/databases/mongodb"
	mongodbpaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/mongodb/interfaces"
	databasepaymentdatasources "github.com/horeekaa/backend/features/payments/data/dataSources/databases/sources"
)

type PaymentDataSourceDependency struct{}

func (_ *PaymentDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbpaymentdatasourceinterfaces.PaymentDataSourceMongo {
			paymentDataSourceMongo, _ := mongodbpaymentdatasources.NewPaymentDataSourceMongo(basicOperation)
			return paymentDataSourceMongo
		},
	)

	container.Singleton(
		func(paymentDataSourceMongo mongodbpaymentdatasourceinterfaces.PaymentDataSourceMongo) databasepaymentdatasourceinterfaces.PaymentDataSource {
			paymentDataSource, _ := databasepaymentdatasources.NewPaymentDataSource()
			paymentDataSource.SetMongoDataSource(paymentDataSourceMongo)
			return paymentDataSource
		},
	)
}
