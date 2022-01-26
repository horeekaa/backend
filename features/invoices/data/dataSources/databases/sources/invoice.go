package databaseinvoicedatasources

import (
	databaseinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/interfaces/sources"
	mongodbinvoicedatasourceinterfaces "github.com/horeekaa/backend/features/invoices/data/dataSources/databases/mongodb/interfaces"
)

type invoiceDataSource struct {
	invoiceDataSourceRepoMongo mongodbinvoicedatasourceinterfaces.InvoiceDataSourceMongo
}

func (invoiceDataSource *invoiceDataSource) SetMongoDataSource(mongoDataSource mongodbinvoicedatasourceinterfaces.InvoiceDataSourceMongo) bool {
	invoiceDataSource.invoiceDataSourceRepoMongo = mongoDataSource
	return true
}

func (invoiceDataSource *invoiceDataSource) GetMongoDataSource() mongodbinvoicedatasourceinterfaces.InvoiceDataSourceMongo {
	return invoiceDataSource.invoiceDataSourceRepoMongo
}

func NewInvoiceDataSource() (databaseinvoicedatasourceinterfaces.InvoiceDataSource, error) {
	return &invoiceDataSource{}, nil
}
