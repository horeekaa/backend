package authenticationcoredependencies

import (
	"github.com/golobby/container/v2"
	authenticationcoreclients "github.com/horeekaa/backend/core/authentication"
	firebaseauthcoreclients "github.com/horeekaa/backend/core/authentication/firebase"
	firebaseauthcoreclientinterfaces "github.com/horeekaa/backend/core/authentication/firebase/interfaces"
	authenticationcoreclientinterfaces "github.com/horeekaa/backend/core/authentication/interfaces/init"
	firebaseserverlesscoreclientinterfaces "github.com/horeekaa/backend/core/serverless/firebase/interfaces"
)

type FirebaseAuthenticationDependency struct{}

func (firebaseAuthDependency *FirebaseAuthenticationDependency) bind() {
	container.Singleton(
		func(fbSvless firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient) firebaseauthcoreclientinterfaces.FirebaseAuthenticationClient {
			fbAuthClient, _ := firebaseauthcoreclients.NewFirebaseAuthClient()
			fbAuthClient.InitializeClient(fbSvless)
			return fbAuthClient
		},
	)

	container.Singleton(
		func(fbAuthClient firebaseauthcoreclientinterfaces.FirebaseAuthenticationClient) authenticationcoreclientinterfaces.AuthenticationClient {
			authClient, _ := authenticationcoreclients.NewAuthClient()
			authClient.SetFirebaseAuthClient(fbAuthClient)
			return authClient
		},
	)
}
