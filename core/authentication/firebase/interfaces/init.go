package firebaseauthcoreclientinterfaces

import (
	auth "firebase.google.com/go/v4/auth"
	firebaseserverlesscoreclientinterfaces "github.com/horeekaa/backend/core/serverless/firebase/interfaces"
)

type FirebaseAuthenticationClient interface {
	InitializeClient(fbServerlessClient firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient) (bool, error)
	GetServerlessClient() (firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient, error)
	GetAuthClient() (*auth.Client, error)
}
