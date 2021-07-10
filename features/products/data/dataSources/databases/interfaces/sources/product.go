package databaseproductdatasourceinterfaces

import (
	mongodbproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/mongodb/interfaces"
)

type ProductDataSource interface {
	GetMongoDataSource() mongodbproductdatasourceinterfaces.ProductDataSourceMongo
	SetMongoDataSource(mongoRepo mongodbproductdatasourceinterfaces.ProductDataSourceMongo) bool
}
