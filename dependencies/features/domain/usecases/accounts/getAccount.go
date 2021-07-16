package accountpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountpresentationusecases "github.com/horeekaa/backend/features/accounts/domain/usecases"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
)

type GetAccountUsecaseDependency struct{}

func (_ GetAccountUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepository accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountRepository accountdomainrepositoryinterfaces.GetAccountRepository,
			getAccountMemberAccessRepository memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
		) accountpresentationusecaseinterfaces.GetAccountUsecase {
			getAccountUsecase, _ := accountpresentationusecases.NewGetAccountUsecase(
				getAccountFromAuthDataRepository,
				getAccountRepository,
				getAccountMemberAccessRepository,
			)
			return getAccountUsecase
		},
	)
}
