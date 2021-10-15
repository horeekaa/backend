package addressregiongrouppresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
	addressregiongrouppresentationusecases "github.com/horeekaa/backend/features/addressRegionGroups/domain/usecases"
	addressregiongrouppresentationusecaseinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/presentation/usecases"
)

type GetAddressRegionGroupUsecaseDependency struct{}

func (_ GetAddressRegionGroupUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAddressRegionGroupRepo addressregiongroupdomainrepositoryinterfaces.GetAddressRegionGroupRepository,
		) addressregiongrouppresentationusecaseinterfaces.GetAddressRegionGroupUsecase {
			getAddressRegionGroupUsecase, _ := addressregiongrouppresentationusecases.NewGetAddressRegionGroupUsecase(
				getAddressRegionGroupRepo,
			)
			return getAddressRegionGroupUsecase
		},
	)
}
