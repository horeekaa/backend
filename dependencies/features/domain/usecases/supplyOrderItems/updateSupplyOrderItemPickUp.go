package supplyorderitempresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	supplyorderitempresentationusecases "github.com/horeekaa/backend/features/supplyOrderItems/domain/usecases"
	supplyorderitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/presentation/usecases"
)

type UpdateSupplyOrderItemPickUpUsecaseDependency struct{}

func (_ *UpdateSupplyOrderItemPickUpUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			proposeUpdateSupplyOrderItemPickUpRepo supplyorderitemdomainrepositoryinterfaces.ProposeUpdateSupplyOrderItemPickUpRepository,
		) supplyorderitempresentationusecaseinterfaces.UpdateSupplyOrderItemPickUpUsecase {
			updateSupplyOrderItemUsecase, _ := supplyorderitempresentationusecases.NewUpdateSupplyOrderItemPickUpUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				proposeUpdateSupplyOrderItemPickUpRepo,
			)
			return updateSupplyOrderItemUsecase
		},
	)
}
