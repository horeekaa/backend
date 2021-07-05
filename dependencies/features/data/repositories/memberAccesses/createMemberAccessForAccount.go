package memberaccessdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	memberaccessdomainrepositories "github.com/horeekaa/backend/features/memberAccesses/data/repositories"
	memberaccessdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccesses/domain/repositories"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
)

type CreateMemberAccessForAccountDependency struct{}

func (createMemberAccessForAccountDependency *CreateMemberAccessForAccountDependency) Bind() {
	container.Singleton(
		func(
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
			memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
			memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
			organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
		) memberaccessdomainrepositoryinterfaces.CreateMemberAccessForAccountRepository {
			createMemberAccessForAccountRepo, _ := memberaccessdomainrepositories.NewCreateMemberAccessForAccountRepository(
				accountDataSource,
				memberAccessDataSource,
				memberAccessRefDataSource,
				organizationDataSource,
			)
			return createMemberAccessForAccountRepo
		},
	)
}
