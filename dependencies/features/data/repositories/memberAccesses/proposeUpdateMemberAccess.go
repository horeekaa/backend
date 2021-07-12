package memberaccessdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositories "github.com/horeekaa/backend/features/memberAccesses/data/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
)

type ProposeUpdateMemberAccessDependency struct{}

func (_ *ProposeUpdateMemberAccessDependency) Bind() {
	container.Singleton(
		func(
			memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
			loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
			memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
			mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
			structComparisonUtility coreutilityinterfaces.StructComparisonUtility,
		) memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessTransactionComponent {
			proposeUpdateMemberAccessComponent, _ := memberaccessdomainrepositories.NewProposeUpdateMemberAccessTransactionComponent(
				memberAccessDataSource,
				loggingDataSource,
				organizationDataSource,
				memberAccessRefDataSource,
				mapProcessorUtility,
				structComparisonUtility,
			)
			return proposeUpdateMemberAccessComponent
		},
	)

	container.Transient(
		func(
			trxComponent memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessTransactionComponent,
			mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
		) memberaccessdomainrepositoryinterfaces.ProposeUpdateMemberAccessRepository {
			proposeUpdateMemberAccessRepo, _ := memberaccessdomainrepositories.NewProposeUpdateMemberAccessRepository(
				trxComponent,
				mongoDBTransaction,
			)
			return proposeUpdateMemberAccessRepo
		},
	)
}