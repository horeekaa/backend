package coredependencies

import (
	authenticationcoredependencies "github.com/horeekaa/backend/dependencies/core/authentication"
	databaseclientdependencies "github.com/horeekaa/backend/dependencies/core/databaseClient"
	serverlesscoredependencies "github.com/horeekaa/backend/dependencies/core/serverless"
	coreutilitydependencies "github.com/horeekaa/backend/dependencies/core/utilities"
	dependencybindinginterfaces "github.com/horeekaa/backend/dependencies/interfaces"
)

type CoreDependency struct{}

func (_ *CoreDependency) Bind() {
	registrationList := [...]dependencybindinginterfaces.BindingInterface{
		&authenticationcoredependencies.FirebaseAuthenticationDependency{},
		&databaseclientdependencies.DatabaseDependency{},
		&serverlesscoredependencies.FirebaseServerlessDependency{},
		&coreutilitydependencies.StructComparisonDependency{},
		&coreutilitydependencies.StructFieldIteratorDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
