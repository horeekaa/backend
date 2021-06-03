package masterdependencies

import (
	coredependencies "github.com/horeekaa/backend/dependencies/core/bindings"
	accountdependencies "github.com/horeekaa/backend/dependencies/features/accounts/bindings"
	loggingdependencies "github.com/horeekaa/backend/dependencies/features/loggings/bindings"
	memberaccessrefdependencies "github.com/horeekaa/backend/dependencies/features/memberAccessRefs/bindings"
	memberaccessdependencies "github.com/horeekaa/backend/dependencies/features/memberAccesses/bindings"
	organizationdependencies "github.com/horeekaa/backend/dependencies/features/organizations/bindings"
	dependencybindinginterfaces "github.com/horeekaa/backend/dependencies/interfaces"
)

type MasterDependency struct{}

func (_ *MasterDependency) Bind() {
	registrationList := []dependencybindinginterfaces.BindingInterface{
		&coredependencies.CoreDependency{},

		&memberaccessdependencies.MemberAccessDependency{},
		&accountdependencies.AccountDependency{},
		&loggingdependencies.LoggingDependency{},
		&memberaccessrefdependencies.MemberAccessRefDependency{},
		&organizationdependencies.OrganizationDependency{},
		&memberaccessdependencies.MemberAccessDependency{},
		&accountdependencies.AccountDependency{},
		&organizationdependencies.OrganizationDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
