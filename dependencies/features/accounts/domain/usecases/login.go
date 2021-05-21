package accountpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountpresentationusecases "github.com/horeekaa/backend/features/accounts/domain/usecases"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
)

type LoginUsecaseDependency struct{}

func (loginUsecaseDpdcy *LoginUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			createAccountFromAuthDataRepo accountdomainrepositoryinterfaces.CreateAccountFromAuthDataRepository,
			createMemberAccessForAccountRepo memberaccessdomainrepositoryinterfaces.CreateMemberAccessForAccountRepository,
			getAccountMemberAccessRepository memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			manageAccountDeviceTokenRepository accountdomainrepositoryinterfaces.ManageAccountDeviceTokenRepository,
		) accountpresentationusecaseinterfaces.LoginUsecase {
			loginUsecase, _ := accountpresentationusecases.NewLoginUsecase(
				getAccountFromAuthDataRepo,
				createAccountFromAuthDataRepo,
				createMemberAccessForAccountRepo,
				getAccountMemberAccessRepository,
				manageAccountDeviceTokenRepository,
			)
			return loginUsecase
		},
	)
}
