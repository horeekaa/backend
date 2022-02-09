package mongodbinvoicedatasourcedependencies

import (
	container "github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databaseinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/interfaces/sources"
	mongodbinvoicedatasources "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/mongodb"
	mongodbinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/mongodb/interfaces"
	databaseinvoicedatasources "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/sources"
)

type InvoiceDataSourceDependency struct{}

func (_ *InvoiceDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbinvoicedatasourceinterfaces.InvoiceDataSourceMongo {
			invoiceDataSourceMongo, _ := mongodbinvoicedatasources.NewInvoiceDataSourceMongo(basicOperation)
			return invoiceDataSourceMongo
		},
	)

	container.Singleton(
		func(invoiceDataSourceMongo mongodbinvoicedatasourceinterfaces.InvoiceDataSourceMongo) databaseinvoicedatasourceinterfaces.InvoiceDataSource {
			invoiceDataSource, _ := databaseinvoicedatasources.NewInvoiceDataSource()
			invoiceDataSource.SetMongoDataSource(invoiceDataSourceMongo)
			return invoiceDataSource
		},
	)
}
