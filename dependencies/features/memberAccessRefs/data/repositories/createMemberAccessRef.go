package memberaccessrefdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	memberaccessrefdomainrepositories "github.com/horeekaa/backend/features/memberAccessRefs/data/repositories"
	memberaccessrefdomainrepositoryinterfaces "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories"
	databasememberaccessrefdatasourceinterfaces "github.com/horeekaa/backend/features/memberaccessrefs/data/dataSources/databases/interfaces/sources"
)

type CreateMemberAccessRefDependency struct{}

func (_ *CreateMemberAccessRefDependency) bind() {
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
