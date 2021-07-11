package productpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	productpresentationusecases "github.com/horeekaa/backend/features/products/domain/usecases"
	productpresentationusecaseinterfaces "github.com/horeekaa/backend/features/products/presentation/usecases"
)

type CreateProductUsecaseDependency struct{}

func (_ *CreateProductUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			createProductRepo productdomainrepositoryinterfaces.CreateProductRepository,
		) productpresentationusecaseinterfaces.CreateProductUsecase {
			productRefUcase, _ := productpresentationusecases.NewCreateProductUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				createProductRepo,
			)
			return productRefUcase
		},
	)
}
