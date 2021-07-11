package productdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	productdomainrepositories "github.com/horeekaa/backend/features/products/data/repositories"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
)

type GetProductDependency struct{}

func (_ *GetProductDependency) Bind() {
	container.Singleton(
		func(
			productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
		) productdomainrepositoryinterfaces.GetProductRepository {
			getproductRepo, _ := productdomainrepositories.NewGetProductRepository(
				productDataSource,
			)
			return getproductRepo
		},
	)
}
