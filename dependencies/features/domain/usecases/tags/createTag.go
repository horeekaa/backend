package tagpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	tagpresentationusecases "github.com/horeekaa/backend/features/tags/domain/usecases"
	tagpresentationusecaseinterfaces "github.com/horeekaa/backend/features/tags/presentation/usecases"
)

type CreateTagUsecaseDependency struct{}

func (_ *CreateTagUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			createTagRepo tagdomainrepositoryinterfaces.CreateTagRepository,
		) tagpresentationusecaseinterfaces.CreateTagUsecase {
			tagUcase, _ := tagpresentationusecases.NewCreateTagUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				createTagRepo,
			)
			return tagUcase
		},
	)
}
