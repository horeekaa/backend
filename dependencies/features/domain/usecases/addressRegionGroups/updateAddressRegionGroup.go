package addressregiongrouppresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
	addressregiongrouppresentationusecases "github.com/horeekaa/backend/features/addressRegionGroups/domain/usecases"
	addressregiongrouppresentationusecaseinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/presentation/usecases"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
)

type UpdateAddressRegionGroupUsecaseDependency struct{}

func (_ *UpdateAddressRegionGroupUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			proposeUpdateaddressRegionGroupRepo addressregiongroupdomainrepositoryinterfaces.ProposeUpdateAddressRegionGroupRepository,
			approveUpdateaddressRegionGroupRepo addressregiongroupdomainrepositoryinterfaces.ApproveUpdateAddressRegionGroupRepository,
		) addressregiongrouppresentationusecaseinterfaces.UpdateAddressRegionGroupUsecase {
			updateAddressRegionGroupUsecase, _ := addressregiongrouppresentationusecases.NewUpdateAddressRegionGroupUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				proposeUpdateaddressRegionGroupRepo,
				approveUpdateaddressRegionGroupRepo,
			)
			return updateAddressRegionGroupUsecase
		},
	)
}
