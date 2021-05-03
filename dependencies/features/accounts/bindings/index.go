package accountdependencies

import (
	firebaseauthdependencies "github.com/horeekaa/backend/dependencies/features/accounts/data/dataSources/authentication"
	mongodbaccountdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/accounts/data/dataSources/databases/mongodb"
	accountdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/accounts/data/repositories"
	accountpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/accounts/domain/usecases"
	dependencybindinginterfaces "github.com/horeekaa/backend/dependencies/interfaces"
)

type AccountDependency struct{}

func (_ *AccountDependency) Bind() {
	registrationList := [...]dependencybindinginterfaces.BindingInterface{
		&firebaseauthdependencies.FirebaseAuthDependency{},
		&mongodbaccountdatasourcedependencies.AccountDataSourceDependency{},
		&mongodbaccountdatasourcedependencies.MemberAccessDataSourceDependency{},
		&mongodbaccountdatasourcedependencies.PersonDataSourceDependency{},

		&accountdomainrepositorydependencies.CreateMemberAccessForAccountDependency{},
		&accountdomainrepositorydependencies.GetAccountDependency{},
		&accountdomainrepositorydependencies.GetAccountMemberAccessDependency{},
		&accountdomainrepositorydependencies.GetPersonDataFromAccountDependency{},
		&accountdomainrepositorydependencies.ManageAccountAuthenticationDependency{},
		&accountdomainrepositorydependencies.ManageAccountDeviceTokenDependency{},

		&accountpresentationusecasedependencies.GetAccountUsecaseDependency{},
		&accountpresentationusecasedependencies.GetPersonDataFromAccountUsecaseDependency{},
		&accountpresentationusecasedependencies.LoginUsecaseDependency{},
		&accountpresentationusecasedependencies.LogoutUsecaseDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
