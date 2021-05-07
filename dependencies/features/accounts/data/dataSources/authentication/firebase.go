package firebaseauthdependencies

import (
	"github.com/golobby/container/v2"
	firebaseauthcoreclientinterfaces "github.com/horeekaa/backend/core/authentication/firebase/interfaces"
	firebaseauthdatasources "github.com/horeekaa/backend/features/accounts/data/dataSources/authentication"
	firebaseauthdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/authentication/interfaces"
)

type FirebaseAuthDependency struct{}

func (firebaseAuthDependency *FirebaseAuthDependency) Bind() {
	container.Singleton(
		func(fbAuthClient firebaseauthcoreclientinterfaces.FirebaseAuthenticationClient) firebaseauthdatasourceinterfaces.FirebaseAuthRepo {
			fbAuthRepo, _ := firebaseauthdatasources.NewFirebaseAuthRepo(fbAuthClient)
			return fbAuthRepo
		},
	)
}
