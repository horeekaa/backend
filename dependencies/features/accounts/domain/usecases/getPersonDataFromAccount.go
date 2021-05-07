package accountpresentationusecasedependencies

import (
	"github.com/golobby/container/v2"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountpresentationusecases "github.com/horeekaa/backend/features/accounts/domain/usecases"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
)

type GetPersonDataFromAccountUsecaseDependency struct{}

func (getPrsnDataFromAccUsecaseDpdcy *GetPersonDataFromAccountUsecaseDependency) Bind() {
	container.Singleton(
		func(
			getAccountFromAuthDataRepo accountdomainrepositoryinterfaces.GetAccountFromAuthData,
			getAccountMemberAccessRepository accountdomainrepositoryinterfaces.GetAccountMemberAccessRepository,
			getPersonDataFromAccountRepository accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository,
		) accountpresentationusecaseinterfaces.GetPersonDataFromAccountUsecase {
			getPersonDataFromAccountUsecase, _ := accountpresentationusecases.NewGetPersonDataFromAccountUsecase(
				getAccountFromAuthDataRepo,
				getAccountMemberAccessRepository,
				getPersonDataFromAccountRepository,
			)
			return getPersonDataFromAccountUsecase
		},
	)
}
