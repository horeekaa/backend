package memberaccessdependencies

import (
	mongodbmemberaccessdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/memberAccesses/data/dataSources/databases/mongodb"
	memberaccessdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/memberAccesses/data/repositories"
	dependencybindinginterfaces "github.com/horeekaa/backend/dependencies/interfaces"
)

type MemberAccessDependency struct{}

func (_ *MemberAccessDependency) Bind() {
	registrationList := [...]dependencybindinginterfaces.BindingInterface{
		&mongodbmemberaccessdatasourcedependencies.MemberAccessDataSourceDependency{},

		&memberaccessdomainrepositorydependencies.CreateMemberAccessForAccountDependency{},
		&memberaccessdomainrepositorydependencies.GetAccountMemberAccessDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
