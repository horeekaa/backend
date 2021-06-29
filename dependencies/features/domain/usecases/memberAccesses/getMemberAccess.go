package memberaccesspresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccesspresentationusecases "github.com/horeekaa/backend/features/memberAccesses/domain/usecases"
	memberaccesspresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases"
)

type GetMemberAccessUsecaseDependency struct{}

func (_ GetMemberAccessUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
		) memberaccesspresentationusecaseinterfaces.GetMemberAccessUsecase {
			getMemberAccessUsecase, _ := memberaccesspresentationusecases.NewGetMemberAccessUsecase(
				getMemberAccessRepo,
			)
			return getMemberAccessUsecase
		},
	)
}
