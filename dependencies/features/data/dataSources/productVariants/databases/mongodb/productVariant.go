package mongodbproductvariantdatasourcedependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbproductvariantdatasources "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/mongodb"
	mongodbproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/mongodb/interfaces"
	databaseproductvariantdatasources "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/sources"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productvariants/data/dataSources/databases/interfaces/sources"
)

type ProductVariantDataSourceDependency struct{}

func (_ *ProductVariantDataSourceDependency) Bind() {
	container.Singleton(
		func(basicOperation mongodbcoreoperationinterfaces.BasicOperation) mongodbproductvariantdatasourceinterfaces.ProductVariantDataSourceMongo {
			productVariantDataSourceMongo, _ := mongodbproductvariantdatasources.NewProductVariantDataSourceMongo(basicOperation)
			return productVariantDataSourceMongo
		},
	)

	container.Singleton(
		func(productVariantDataSourceMongo mongodbproductvariantdatasourceinterfaces.ProductVariantDataSourceMongo) databaseproductvariantdatasourceinterfaces.ProductVariantDataSource {
			productVariantRepo, _ := databaseproductvariantdatasources.NewProductVariantDataSource()
			productVariantRepo.SetMongoDataSource(productVariantDataSourceMongo)
			return productVariantRepo
		},
	)
}
