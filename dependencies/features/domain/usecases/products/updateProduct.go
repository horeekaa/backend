package productpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	productpresentationusecases "github.com/horeekaa/backend/features/products/domain/usecases"
	productpresentationusecaseinterfaces "github.com/horeekaa/backend/features/products/presentation/usecases"
)

type UpdateProductUsecaseDependency struct{}

func (_ *UpdateProductUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			proposeUpdateproductRepo productdomainrepositoryinterfaces.ProposeUpdateProductRepository,
			approveUpdateproductRepo productdomainrepositoryinterfaces.ApproveUpdateProductRepository,
		) productpresentationusecaseinterfaces.UpdateProductUsecase {
			updateProductUsecase, _ := productpresentationusecases.NewUpdateProductUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				proposeUpdateproductRepo,
				approveUpdateproductRepo,
			)
			return updateProductUsecase
		},
	)
}
