package accountdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	firebaseauthdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/authentication/interfaces"
	accountdomainrepositories "github.com/horeekaa/backend/features/accounts/data/repositories"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
)

type GetUserFromAuthHeaderDependency struct{}

func (_ *GetUserFromAuthHeaderDependency) Bind() {
	container.Singleton(
		func(
			firebaseDataSource firebaseauthdatasourceinterfaces.FirebaseAuthRepo,
		) accountdomainrepositoryinterfaces.GetUserFromAuthHeaderRepository {
			getUserFromAuthHeaderRepo, _ := accountdomainrepositories.NewGetUserFromAuthHeaderRepository(
				firebaseDataSource,
			)
			return getUserFromAuthHeaderRepo
		},
	)
}
