package productvariantpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	productvariantdomainrepositoryinterfaces "github.com/horeekaa/backend/features/productVariants/domain/repositories"
	productvariantpresentationusecases "github.com/horeekaa/backend/features/productVariants/domain/usecases"
	productvariantpresentationusecaseinterfaces "github.com/horeekaa/backend/features/productVariants/presentation/usecases"
)

type GetProductVariantUsecaseDependency struct{}

func (_ GetProductVariantUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getProductVariantRepo productvariantdomainrepositoryinterfaces.GetProductVariantRepository,
		) productvariantpresentationusecaseinterfaces.GetProductVariantUsecase {
			getProductVariantUsecase, _ := productvariantpresentationusecases.NewGetProductVariantUsecase(
				getProductVariantRepo,
			)
			return getProductVariantUsecase
		},
	)
}
