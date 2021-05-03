package firebaseauthcoreclients

import (
	"context"
	"errors"

	auth "firebase.google.com/go/v4/auth"

	firebaseauthcoreclientinterfaces "github.com/horeekaa/backend/core/authentication/firebase/interfaces"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/repoExceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/repoExceptions/enums"
	firebaseserverlesscoreclientinterfaces "github.com/horeekaa/backend/core/serverless/firebase/interfaces"
)

type firebaseAuthenticationClient struct {
	firebaseServerlessClient firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient
	client                   *auth.Client
}

func (fbAuthClient *firebaseAuthenticationClient) GetServerlessClient() (firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient, error) {
	if fbAuthClient.firebaseServerlessClient == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClientInitializationFailed,
			"/newFirebaseAuthentication",
			errors.New(horeekaacoreexceptionenums.ClientInitializationFailed),
		)
	}
	return fbAuthClient.firebaseServerlessClient, nil
}

func (fbAuthClient *firebaseAuthenticationClient) GetAuthClient() (*auth.Client, error) {
	if fbAuthClient.client == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClientInitializationFailed,
			"/newFirebaseAuthentication",
			errors.New(horeekaacoreexceptionenums.ClientInitializationFailed),
		)
	}
	return fbAuthClient.client, nil
}

func (fbAuthClient *firebaseAuthenticationClient) InitializeClient(firebaseClient firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient) (bool, error) {
	app, _ := firebaseClient.GetApp()
	client, err := (*app).Auth(context.Background())
	if err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			"/newFirebaseAuthentication",
			err,
		)
	}
	fbAuthClient.firebaseServerlessClient = firebaseClient
	fbAuthClient.client = client

	return true, nil
}

func NewFirebaseAuthClient() (firebaseauthcoreclientinterfaces.FirebaseAuthenticationClient, error) {
	return &firebaseAuthenticationClient{}, nil
}
