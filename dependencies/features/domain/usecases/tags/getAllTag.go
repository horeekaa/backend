package tagpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	tagpresentationusecases "github.com/horeekaa/backend/features/tags/domain/usecases"
	tagpresentationusecaseinterfaces "github.com/horeekaa/backend/features/tags/presentation/usecases"
)

type GetAllTagUsecaseDependency struct{}

func (_ GetAllTagUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getAllTagRepo tagdomainrepositoryinterfaces.GetAllTagRepository,
		) tagpresentationusecaseinterfaces.GetAllTagUsecase {
			getAllTagUcase, _ := tagpresentationusecases.NewGetAllTagUsecase(
				getAccountFromAuthDataRepo,
				getAccountRepo,
				getAllTagRepo,
			)
			return getAllTagUcase
		},
	)
}
