package productpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	productpresentationusecases "github.com/horeekaa/backend/features/products/domain/usecases"
	productpresentationusecaseinterfaces "github.com/horeekaa/backend/features/products/presentation/usecases"
)

type GetAllProductUsecaseDependency struct{}

func (_ GetAllProductUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getAllProductRepo productdomainrepositoryinterfaces.GetAllProductRepository,
		) productpresentationusecaseinterfaces.GetAllProductUsecase {
			getAllProductUcase, _ := productpresentationusecases.NewGetAllProductUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				getAllProductRepo,
			)
			return getAllProductUcase
		},
	)
}
