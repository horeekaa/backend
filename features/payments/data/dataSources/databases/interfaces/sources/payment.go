package databasepaymentdatasourceinterfaces

import (
	mongodbpaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/mongodb/interfaces"
)

type PaymentDataSource interface {
	GetMongoDataSource() mongodbpaymentdatasourceinterfaces.PaymentDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbpaymentdatasourceinterfaces.PaymentDataSourceMongo) bool
}
