package firebaseauthcoreclientinterfaces

import (
	firebaseauthcorewrapperinterfaces "github.com/horeekaa/backend/core/authentication/firebase/interfaces/wrappers"
	firebaseserverlesscoreclientinterfaces "github.com/horeekaa/backend/core/serverless/firebase/interfaces"
)

type FirebaseAuthenticationClient interface {
	InitializeClient(fbServerlessClient firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient) (bool, error)
	GetServerlessClient() (firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient, error)
	GetAuthClient() (firebaseauthcorewrapperinterfaces.FirebaseAuthClient, error)
}
