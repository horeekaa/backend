package coredependencies

import (
	authenticationcoredependencies "github.com/horeekaa/backend/dependencies/core/authentication"
	databaseclientdependencies "github.com/horeekaa/backend/dependencies/core/databaseClient"
	i18ncoredependencies "github.com/horeekaa/backend/dependencies/core/i18n"
	messagingdependencies "github.com/horeekaa/backend/dependencies/core/messaging"
	serverlesscoredependencies "github.com/horeekaa/backend/dependencies/core/serverless"
	storagecoredependencies "github.com/horeekaa/backend/dependencies/core/storages"
	coreutilitydependencies "github.com/horeekaa/backend/dependencies/core/utilities"
	dependencybindinginterfaces "github.com/horeekaa/backend/dependencies/interfaces"
)

type CoreDependency struct{}

func (_ *CoreDependency) Bind() {
	registrationList := [...]dependencybindinginterfaces.BindingInterface{
		&coreutilitydependencies.MapProcessorUtilityDependency{},
		&serverlesscoredependencies.FirebaseServerlessDependency{},
		&authenticationcoredependencies.FirebaseAuthenticationDependency{},
		&messagingdependencies.FirebaseMessagingDependency{},
		&databaseclientdependencies.DatabaseDependency{},
		&storagecoredependencies.GoogleCloudStorageDependency{},
		&i18ncoredependencies.GoLocalizeI18NDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
