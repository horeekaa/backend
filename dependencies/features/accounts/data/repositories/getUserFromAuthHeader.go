package accountdomainrepositorydependencies

import (
	"github.com/golobby/container/v2"
	authenticationcoreclientinterfaces "github.com/horeekaa/backend/core/authentication/interfaces"
	accountdomainrepositories "github.com/horeekaa/backend/features/accounts/data/repositories"
	accountdomainrepositoryinterfaces "github.com/horeekaa/backend/features/accounts/domain/repositories"
)

type GetUserFromAuthHeaderDependency struct{}

func (_ *GetUserFromAuthHeaderDependency) Bind() {
	container.Singleton(
		func(
			firebaseDataSource authenticationcoreclientinterfaces.AuthenticationRepo,
		) accountdomainrepositoryinterfaces.GetUserFromAuthHeaderRepository {
			getUserFromAuthHeaderRepo, _ := accountdomainrepositories.NewGetUserFromAuthHeaderRepository(
				firebaseDataSource,
			)
			return getUserFromAuthHeaderRepo
		},
	)
}
