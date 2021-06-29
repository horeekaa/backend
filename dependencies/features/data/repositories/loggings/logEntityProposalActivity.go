package loggingdomainrepositorydependencies

import (
	container "github.com/golobby/container/v2"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	loggingdomainrepositories "github.com/horeekaa/backend/features/loggings/data/repositories"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
)

type LogEntityProposalActivityDependency struct{}

func (logEntityProposalActivityDependency *LogEntityProposalActivityDependency) Bind() {
	container.Singleton(
		func(
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			structComparisonUtility coreutilityinterfaces.StructComparisonUtility,
			structFieldIteratorUtility coreutilityinterfaces.StructFieldIteratorUtility,
		) loggingdomainrepositoryinterfaces.LogEntityProposalActivityRepository {
			logEntityProposalActivityRepo, _ := loggingdomainrepositories.NewLogEntityProposalActivityRepository(
				loggingDataSource,
				structComparisonUtility,
				structFieldIteratorUtility,
			)
			return logEntityProposalActivityRepo
		},
	)
}
