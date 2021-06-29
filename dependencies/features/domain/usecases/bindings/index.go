package usecasesdependencies

import (
	accountpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/accounts"
	loggingpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/loggings"
	memberaccessrefpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/memberAccessRefs"
	memberaccesspresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/memberAccesses"
	organizationpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/organizations"
	dependencybindinginterfaces "github.com/horeekaa/backend/dependencies/interfaces"
)

type UsecasesDependency struct{}

func (_ *UsecasesDependency) Bind() {
	registrationList := [...]dependencybindinginterfaces.BindingInterface{
		&accountpresentationusecasedependencies.GetAccountUsecaseDependency{},
		&accountpresentationusecasedependencies.GetPersonDataFromAccountUsecaseDependency{},
		&accountpresentationusecasedependencies.LoginUsecaseDependency{},
		&accountpresentationusecasedependencies.LogoutUsecaseDependency{},
		&accountpresentationusecasedependencies.GetAuthUserAndAttachToCtxUsecaseDependency{},

		&loggingpresentationusecasedependencies.GetLoggingUsecaseDependency{},

		&memberaccesspresentationusecasedependencies.CreateMemberAccessUsecaseDependency{},
		&memberaccesspresentationusecasedependencies.GetAllMemberAccessUsecaseDependency{},
		&memberaccesspresentationusecasedependencies.GetMemberAccessUsecaseDependency{},
		&memberaccesspresentationusecasedependencies.UpdateMemberAccessUsecaseDependency{},

		&memberaccessrefpresentationusecasedependencies.CreateMemberAccessRefUsecaseDependency{},
		&memberaccessrefpresentationusecasedependencies.GetAllMemberAccessRefUsecaseDependency{},
		&memberaccessrefpresentationusecasedependencies.GetMemberAccessRefUsecaseDependency{},
		&memberaccessrefpresentationusecasedependencies.UpdateMemberAccessRefUsecaseDependency{},

		&organizationpresentationusecasedependencies.CreateOrganizationUsecaseDependency{},
		&organizationpresentationusecasedependencies.GetAllOrganizationUsecaseDependency{},
		&organizationpresentationusecasedependencies.GetOrganizationUsecaseDependency{},
		&organizationpresentationusecasedependencies.UpdateOrganizationUsecaseDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
