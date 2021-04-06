package accountdomainusecasedependencies

import (
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountpresentationusecases "github.com/horeekaa/backend/features/accounts/domain/usecases"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
)

type LoginUsecaseDependency struct{}

func (loginUsecaseDpdcy *LoginUsecaseDependency) bind(
	manageAccountAuthenticationRepository accountdomainrepositoryinterfaces.ManageAccountAuthenticationRepository,
	manageAccountDeviceTokenRepository accountdomainrepositoryinterfaces.ManageAccountDeviceTokenRepository,
) accountpresentationusecaseinterfaces.LoginUsecase {
	loginUsecase, _ := accountpresentationusecases.NewLoginUsecase(
		manageAccountAuthenticationRepository,
		manageAccountDeviceTokenRepository,
	)
	return loginUsecase
}
