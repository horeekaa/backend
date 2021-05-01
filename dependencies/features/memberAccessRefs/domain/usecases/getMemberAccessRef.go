package memberaccessrefpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	memberaccessrefpresentationusecases "github.com/horeekaa/backend/features/memberAccessRefs/domain/usecases"
	memberaccessrefpresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases"
)

type GetMemberAccessRefUsecaseDependency struct{}

func (_ GetMemberAccessRefUsecaseDependency) bind() {
	container.Singleton(
		func(
			getMemberAccessRefRepo memberaccessrefdomainrepositoryinterfaces.GetMemberAccessRefRepository,
		) memberaccessrefpresentationusecaseinterfaces.GetMemberAccessRefUsecase {
			getMemberAccessRefUsecase, _ := memberaccessrefpresentationusecases.NewGetMemberAccessRefUsecase(
				getMemberAccessRefRepo,
			)
			return getMemberAccessRefUsecase
		},
	)
}
