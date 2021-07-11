package productpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	productpresentationusecases "github.com/horeekaa/backend/features/products/domain/usecases"
	productpresentationusecaseinterfaces "github.com/horeekaa/backend/features/products/presentation/usecases"
)

type GetProductUsecaseDependency struct{}

func (_ GetProductUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getProductRepo productdomainrepositoryinterfaces.GetProductRepository,
		) productpresentationusecaseinterfaces.GetProductUsecase {
			getProductUsecase, _ := productpresentationusecases.NewGetProductUsecase(
				getProductRepo,
			)
			return getProductUsecase
		},
	)
}
