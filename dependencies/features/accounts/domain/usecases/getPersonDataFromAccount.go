package accountdomainusecasedependencies

import (
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	accountpresentationusecases "github.com/horeekaa/backend/features/accounts/domain/usecases"
	accountpresentationusecaseinterfaces "github.com/horeekaa/backend/features/accounts/presentation/usecases"
)

type GetPersonDataFromAccountUsecaseDependency struct{}

func (getPrsnDataFromAccUsecaseDpdcy *GetPersonDataFromAccountUsecaseDependency) bind(
	getPersonDataFromAccountRepository accountdomainrepositoryinterfaces.GetPersonDataFromAccountRepository,
) accountpresentationusecaseinterfaces.GetPersonDataFromAccountUsecase {
	getPersonDataFromAccountUsecase, _ := accountpresentationusecases.NewGetPersonDataFromAccountUsecase(
		getPersonDataFromAccountRepository,
	)
	return getPersonDataFromAccountUsecase
}
