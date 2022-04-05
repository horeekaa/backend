package accountpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountpresentationusecases "github.com/horeekaa/backend/features/accounts/domain/usecases"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
)

type UpdateAccountUsecaseDependency struct{}

func (_ *UpdateAccountUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepo memberaccessdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			updateAccountRepo accountdomainrepositoryinterfaces.UpdateAccountRepository,
		) accountpresentationusecaseinterfaces.UpdateAccountUsecase {
			updateAccountUsecase, _ := accountpresentationusecases.NewUpdateAccountUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepo,
				updateAccountRepo,
			)
			return updateAccountUsecase
		},
	)
}
