package loggingdependencies

import (
	mongodbloggingdatasourcedependencies "github.com/horeekaa/backend/dependencies/features/loggings/data/dataSources/databases/mongodb"
	loggingdomainrepositorydependencies "github.com/horeekaa/backend/dependencies/features/loggings/data/repositories"
	loggingpresentationusecasedependencies "github.com/horeekaa/backend/dependencies/features/loggings/domain/usecases"
	dependencybindinginterfaces "github.com/horeekaa/backend/dependencies/interfaces"
)

type LoggingDependency struct{}

func (_ *LoggingDependency) Bind() {
	registrationList := [...]dependencybindinginterfaces.BindingInterface{
		&mongodbloggingdatasourcedependencies.LoggingDataSourceDependency{},

		&loggingdomainrepositorydependencies.GetLoggingDependency{},
		&loggingdomainrepositorydependencies.LogEntityApprovalActivityDependency{},
		&loggingdomainrepositorydependencies.LogEntityProposalActivityDependency{},

		&loggingpresentationusecasedependencies.GetLoggingUsecaseDependency{},
	}

	for _, reg := range registrationList {
		reg.Bind()
	}
}
