package memberaccessrefdependencies

import (
	mongodbmemberaccessrefdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/memberAccessRefs/data/dataSources/databases/mongodb"
	memberaccessrefdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/memberAccessRefs/data/repositories"
	memberaccessrefpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/memberAccessRefs/domain/usecases"
	dependencybindinginterfaces "github.com/horeekaa/backend/dependencies/interfaces"
)

type MemberAccessRefDependency struct{}

func (_ *MemberAccessRefDependency) Bind() {
	registrationList := [...]dependencybindinginterfaces.BindingInterface{
		&mongodbmemberaccessrefdatasourcedependencies.MemberAccessRefDataSourceDependency{},

		&memberaccessrefdomainrepositorydependencies.CreateMemberAccessRefDependency{},
		&memberaccessrefdomainrepositorydependencies.GetAllMemberAccessRefDependency{},
		&memberaccessrefdomainrepositorydependencies.GetMemberAccessRefDependency{},
		&memberaccessrefdomainrepositorydependencies.UpdateMemberAccessRefDependency{},

		&memberaccessrefpresentationusecasedependencies.CreateMemberAccessRefUsecaseDependency{},
		&memberaccessrefpresentationusecasedependencies.GetAllMemberAccessRefUsecaseDependency{},
		&memberaccessrefpresentationusecasedependencies.GetMemberAccessRefUsecaseDependency{},
		&memberaccessrefpresentationusecasedependencies.UpdateMemberAccessRefUsecaseDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
