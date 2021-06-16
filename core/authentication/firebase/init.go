package firebaseauthcoreclients

import (
	"context"

	firebaseauthcoreclientinterfaces "github.com/horeekaa/backend/core/authentication/firebase/interfaces"
	firebaseauthcorewrapperinterfaces "github.com/horeekaa/backend/core/authentication/firebase/interfaces/wrappers"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	firebaseserverlesscoreclientinterfaces "github.com/horeekaa/backend/core/serverless/firebase/interfaces"
)

type firebaseAuthenticationClient struct {
	firebaseServerlessClient firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient
	client                   firebaseauthcorewrapperinterfaces.FirebaseAuthClient
}

func (fbAuthClient *firebaseAuthenticationClient) GetServerlessClient() (firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient, error) {
	if fbAuthClient.firebaseServerlessClient == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClientInitializationFailed,
			"/newFirebaseAuthentication",
			nil,
		)
	}
	return fbAuthClient.firebaseServerlessClient, nil
}

func (fbAuthClient *firebaseAuthenticationClient) GetAuthClient() (firebaseauthcorewrapperinterfaces.FirebaseAuthClient, error) {
	if fbAuthClient.client == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClientInitializationFailed,
			"/newFirebaseAuthentication",
			nil,
		)
	}
	return fbAuthClient.client, nil
}

func (fbAuthClient *firebaseAuthenticationClient) InitializeClient(firebaseClient firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient) (bool, error) {
	app, _ := firebaseClient.GetApp()
	client, err := app.Auth(context.Background())
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
