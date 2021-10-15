package addressregiongrouppresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
	addressregiongrouppresentationusecases "github.com/horeekaa/backend/features/addressRegionGroups/domain/usecases"
	addressregiongrouppresentationusecaseinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/presentation/usecases"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
)

type GetAllAddressRegionGroupUsecaseDependency struct{}

func (_ GetAllAddressRegionGroupUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getAlladdressRegionGroupRepo addressregiongroupdomainrepositoryinterfaces.GetAllAddressRegionGroupRepository,
		) addressregiongrouppresentationusecaseinterfaces.GetAllAddressRegionGroupUsecase {
			getAllAddressRegionGroupUcase, _ := addressregiongrouppresentationusecases.NewGetAllAddressRegionGroupUsecase(
				getAccountFromAuthDataRepo,
				getAccountRepo,
				getAlladdressRegionGroupRepo,
			)
			return getAllAddressRegionGroupUcase
		},
	)
}
