package accountdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	accountdomainrepositories "github.com/horeekaa/backend/features/accounts/data/repositories"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
)

type CreateMemberAccessForAccountDependency struct{}

func (createMemberAccessForAccountDependency *CreateMemberAccessForAccountDependency) Bind() {
	container.Singleton(
		func(
			accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
			memberAccessDataSource databaseaccountdatasourceinterfaces.MemberAccessDataSource,
			memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
		) accountdomainrepositoryinterfaces.CreateMemberAccessForAccountRepository {
			createMemberAccessForAccountRepo, _ := accountdomainrepositories.NewCreateMemberAccessForAccountRepository(
				accountDataSource,
				memberAccessDataSource,
				memberAccessRefDataSource,
			)
			return createMemberAccessForAccountRepo
		},
	)
}
