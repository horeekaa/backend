package masterdependencies

import (
	coredependencies "github.com/horeekaa/backend/dependencies/core/bindings"
	accountdependencies "github.com/horeekaa/backend/dependencies/features/accounts/bindings"
	loggingdependencies "github.com/horeekaa/backend/dependencies/features/loggings/bindings"
	memberaccessrefdependencies "github.com/horeekaa/backend/dependencies/features/memberAccessRefs/bindings"
	organizationdependencies "github.com/horeekaa/backend/dependencies/features/organizations/bindings"
	dependencybindinginterfaces "github.com/horeekaa/backend/dependencies/interfaces"
)

type MasterDependency struct{}

func (_ *MasterDependency) Bind() {
	registrationList := []dependencybindinginterfaces.BindingInterface{
		&coredependencies.CoreDependency{},

		&accountdependencies.AccountDependency{},
		&loggingdependencies.LoggingDependency{},
		&memberaccessrefdependencies.MemberAccessRefDependency{},
		&organizationdependencies.OrganizationDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
