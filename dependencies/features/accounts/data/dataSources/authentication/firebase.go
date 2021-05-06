package firebaseauthdependencies

import (
	"github.com/golobby/container/v2"
	firebaseauthcoreclientinterfaces "github.com/horeekaa/backend/core/authentication/firebase/interfaces"
	authenticationcoreclientinterfaces "github.com/horeekaa/backend/core/authentication/interfaces"
	firebaseauthdatasources "github.com/horeekaa/backend/features/accounts/data/dataSources/authentication"
)

type FirebaseAuthDependency struct{}

func (firebaseAuthDependency *FirebaseAuthDependency) Bind() {
	container.Singleton(
		func(fbAuthClient firebaseauthcoreclientinterfaces.FirebaseAuthenticationClient) authenticationcoreclientinterfaces.AuthenticationRepo {
			fbAuthRepo, _ := firebaseauthdatasources.NewFirebaseAuthRepo(fbAuthClient)
			return fbAuthRepo
		},
	)
}
