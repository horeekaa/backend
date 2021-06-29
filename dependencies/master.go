package masterdependencies

import (
	coredependencies "github.com/horeekaa/backend/dependencies/core/bindings"
	datasourcesdependencies "github.com/horeekaa/backend/dependencies/features/data/dataSources/bindings"
	repositoriesdependencies "github.com/horeekaa/backend/dependencies/features/data/repositories/bindings"
	usecasesdependencies "github.com/horeekaa/backend/dependencies/features/domain/usecases/bindings"
	dependencybindinginterfaces "github.com/horeekaa/backend/dependencies/interfaces"
)

type MasterDependency struct{}

func (_ *MasterDependency) Bind() {
	registrationList := []dependencybindinginterfaces.BindingInterface{
		&coredependencies.CoreDependency{},

		&datasourcesdependencies.DatasourcesDependency{},
		&repositoriesdependencies.RepositoriesDependency{},
		&usecasesdependencies.UsecasesDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
