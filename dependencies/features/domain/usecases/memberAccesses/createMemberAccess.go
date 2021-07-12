package memberaccesspresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	memberaccesspresentationusecases "github.com/horeekaa/backend/features/memberAccesses/domain/usecases"
	memberaccesspresentationusecaseinterfaces "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases"
)

type CreateMemberAccessUsecaseDependency struct{}

func (_ *CreateMemberAccessUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			createMemberAccessRepo memberaccessdomainrepositoryinterfaces.CreateMemberAccessRepository,
		) memberaccesspresentationusecaseinterfaces.CreateMemberAccessUsecase {
			createMmbAccUcase, _ := memberaccesspresentationusecases.NewCreateMemberAccessUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				createMemberAccessRepo,
			)
			return createMmbAccUcase
		},
	)
}
