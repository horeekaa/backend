package authenticationcoreclients

import (
	firebaseauthcoreclientinterfaces "github.com/horeekaa/backend/core/authentication/firebase/interfaces"
	authenticationcoreclientinterfaces "github.com/horeekaa/backend/core/authentication/interfaces"
)

type authClient struct {
	firebaseAuthClient firebaseauthcoreclientinterfaces.FirebaseAuthenticationClient
}

func (athClient *authClient) SetFirebaseAuthClient(firebaseAuthClient firebaseauthcoreclientinterfaces.FirebaseAuthenticationClient) bool {
	athClient.firebaseAuthClient = firebaseAuthClient
	return true
}

func (athClient *authClient) GetFirebaseAuthClient() firebaseauthcoreclientinterfaces.FirebaseAuthenticationClient {
	return athClient.firebaseAuthClient
}

func NewAuthClient() (authenticationcoreclientinterfaces.AuthenticationClient, error) {
	return &authClient{}, nil
}
