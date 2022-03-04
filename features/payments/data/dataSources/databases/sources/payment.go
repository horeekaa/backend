package databasepaymentdatasources

import (
	databasepaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/interfaces/sources"
	mongodbpaymentdatasourceinterfaces "github.com/horeekaa/backend/features/payments/data/dataSources/databases/mongodb/interfaces"
)

type paymentDataSource struct {
	paymentDataSourceRepoMongo mongodbpaymentdatasourceinterfaces.PaymentDataSourceMongo
}

func (paymentDataSource *paymentDataSource) SetMongoDataSource(mongoDataSource mongodbpaymentdatasourceinterfaces.PaymentDataSourceMongo) bool {
	paymentDataSource.paymentDataSourceRepoMongo = mongoDataSource
	return true
}

func (paymentDataSource *paymentDataSource) GetMongoDataSource() mongodbpaymentdatasourceinterfaces.PaymentDataSourceMongo {
	return paymentDataSource.paymentDataSourceRepoMongo
}

func NewPaymentDataSource() (databasepaymentdatasourceinterfaces.PaymentDataSource, error) {
	return &paymentDataSource{}, nil
}
