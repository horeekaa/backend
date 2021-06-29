package loggingdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	loggingdomainrepositories "github.com/horeekaa/backend/features/loggings/data/repositories"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
)

type LogEntityApprovalActivityDependency struct{}

func (logEntityApprActivityDpdcy *LogEntityApprovalActivityDependency) Bind() {
	container.Singleton(
		func(
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
		) loggingdomainrepositoryinterfaces.LogEntityApprovalActivityRepository {
			logEntityApprovalRepo, _ := loggingdomainrepositories.NewLogEntityApprovalActivityRepository(
				loggingDataSource,
			)
			return logEntityApprovalRepo
		},
	)
}
