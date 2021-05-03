package loggingdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	loggingdomainrepositories "github.com/horeekaa/backend/features/loggings/data/repositories"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
)

type GetLoggingDependency struct{}

func (_ *GetLoggingDependency) Bind() {
	container.Singleton(
		func(
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
		) loggingdomainrepositoryinterfaces.GetLoggingRepository {
			getLoggingRepo, _ := loggingdomainrepositories.NewGetLoggingRepository(
				loggingDataSource,
			)
			return getLoggingRepo
		},
	)
}
