package productvariantdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	productvariantdomainrepositories "github.com/horeekaa/backend/features/productVariants/data/repositories"
	productvariantdomainrepositoryinterfaces "github.com/horeekaa/backend/features/productVariants/domain/repositories"
)

type GetProductVariantDependency struct{}

func (_ *GetProductVariantDependency) Bind() {
	container.Singleton(
		func(
			productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource,
		) productvariantdomainrepositoryinterfaces.GetProductVariantRepository {
			getProductVariantRepo, _ := productvariantdomainrepositories.NewGetProductVariantRepository(
				productVariantDataSource,
			)
			return getProductVariantRepo
		},
	)
}
