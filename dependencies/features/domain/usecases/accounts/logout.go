package accountpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountpresentationusecases "github.com/horeekaa/backend/features/accounts/domain/usecases"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
)

type LogoutUsecaseDependency struct{}

func (_ LogoutUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			manageAccountDeviceTokenRepository accountdomainrepositoryinterfaces.ManageAccountDeviceTokenRepository,
		) accountpresentationusecaseinterfaces.LogoutUsecase {
			logoutUsecase, _ := accountpresentationusecases.NewLogoutUsecase(
				getAccountFromAuthDataRepo,
				manageAccountDeviceTokenRepository,
			)
			return logoutUsecase
		},
	)
}
