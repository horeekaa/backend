package databaseproductdatasources

import (
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	mongodbproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/mongodb/interfaces"
)

type productDataSource struct {
	productDataSourceRepoMongo mongodbproductdatasourceinterfaces.ProductDataSourceMongo
}

func (orgDataSource *productDataSource) SetMongoDataSource(mongoDataSource mongodbproductdatasourceinterfaces.ProductDataSourceMongo) bool {
	orgDataSource.productDataSourceRepoMongo = mongoDataSource
	return true
}

func (orgDataSource *productDataSource) GetMongoDataSource() mongodbproductdatasourceinterfaces.ProductDataSourceMongo {
	return orgDataSource.productDataSourceRepoMongo
}

func NewProductDataSource() (databaseproductdatasourceinterfaces.ProductDataSource, error) {
	return &productDataSource{}, nil
}
