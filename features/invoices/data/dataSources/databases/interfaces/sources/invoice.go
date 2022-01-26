package databaseinvoicedatasourceinterfaces

import (
	mongodbinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/mongodb/interfaces"
)

type InvoiceDataSource interface {
	GetMongoDataSource() mongodbinvoicedatasourceinterfaces.InvoiceDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbinvoicedatasourceinterfaces.InvoiceDataSourceMongo) bool
}
