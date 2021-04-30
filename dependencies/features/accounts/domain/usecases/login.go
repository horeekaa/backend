package accountpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountpresentationusecases "github.com/horeekaa/backend/features/accounts/domain/usecases"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
)

type LoginUsecaseDependency struct{}

func (loginUsecaseDpdcy *LoginUsecaseDependency) bind() {
	container.Singleton(
		func(
			manageAccountAuthenticationRepository accountdomainrepositoryinterfaces.ManageAccountAuthenticationRepository,
			getAccountMemberAccessRepository accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			manageAccountDeviceTokenRepository accountdomainrepositoryinterfaces.ManageAccountDeviceTokenRepository,
		) accountpresentationusecaseinterfaces.LoginUsecase {
			loginUsecase, _ := accountpresentationusecases.NewLoginUsecase(
				manageAccountAuthenticationRepository,
				getAccountMemberAccessRepository,
				manageAccountDeviceTokenRepository,
			)
			return loginUsecase
		},
	)
}
