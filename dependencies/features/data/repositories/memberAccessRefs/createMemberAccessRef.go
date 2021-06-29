package memberaccessrefdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/data/dataSources/databases/interfaces/sources"
	memberaccessrefdomainrepositories "github.com/horeekaa/backend/features/memberAccessRefs/data/repositories"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
)

type CreateMemberAccessRefDependency struct{}

func (_ *CreateMemberAccessRefDependency) Bind() {
	container.Singleton(
		func(
			memberAccessRefDataSource databasememberaccessrefdatasourceinterfaces.MemberAccessRefDataSource,
		) memberaccessrefdomainrepositoryinterfaces.CreateMemberAccessRefRepository {
			createMemberAccessRefRepo, _ := memberaccessrefdomainrepositories.NewCreateMemberAccessRefRepository(
				memberAccessRefDataSource,
			)
			return createMemberAccessRefRepo
		},
	)
}
