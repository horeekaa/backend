package mongodbproductdatasourcedependencies

import (
	container "github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	mongodbproductdatasources "github.com/horeekaa/backend/features/products/data/dataSources/databases/mongodb"
	mongodbproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/mongodb/interfaces"
	databaseproductdatasources "github.com/horeekaa/backend/features/products/data/dataSources/databases/sources"
)

type ProductDataSourceDependency struct{}

func (orgDataSourceDpdcy *ProductDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbproductdatasourceinterfaces.ProductDataSourceMongo {
			productDataSourceMongo, _ := mongodbproductdatasources.NewProductDataSourceMongo(basicOperation)
			return productDataSourceMongo
		},
	)

	container.Singleton(
		func(productDataSourceMongo mongodbproductdatasourceinterfaces.ProductDataSourceMongo) databaseproductdatasourceinterfaces.ProductDataSource {
			productRepo, _ := databaseproductdatasources.NewProductDataSource()
			productRepo.SetMongoDataSource(productDataSourceMongo)
			return productRepo
		},
	)
}
