package accountpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountpresentationusecases "github.com/horeekaa/backend/features/accounts/domain/usecases"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
)

type GetPersonDataFromAccountUsecaseDependency struct{}

func (getPrsnDataFromAccUsecaseDpdcy *GetPersonDataFromAccountUsecaseDependency) bind() {
	container.Singleton(
		func(
			manageAccountAuthenticationRepository accountdomainrepositoryinterfaces.ManageAccountAuthenticationRepository,
			getAccountMemberAccessRepository accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getPersonDataFromAccountRepository accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository,
		) accountpresentationusecaseinterfaces.GetPersonDataFromAccountUsecase {
			getPersonDataFromAccountUsecase, _ := accountpresentationusecases.NewGetPersonDataFromAccountUsecase(
				manageAccountAuthenticationRepository,
				getAccountMemberAccessRepository,
				getPersonDataFromAccountRepository,
			)
			return getPersonDataFromAccountUsecase
		},
	)
}
