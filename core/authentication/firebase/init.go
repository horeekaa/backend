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
	pathIdentity             string
}

func (fbAuthClient *firebaseAuthenticationClient) GetServerlessClient() (firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient, error) {
	if fbAuthClient.firebaseServerlessClient == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClientInitializationFailed,
			fbAuthClient.pathIdentity,
			nil,
		)
	}
	return fbAuthClient.firebaseServerlessClient, nil
}

func (fbAuthClient *firebaseAuthenticationClient) GetAuthClient() (firebaseauthcorewrapperinterfaces.FirebaseAuthClient, error) {
	if fbAuthClient.client == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClientInitializationFailed,
			fbAuthClient.pathIdentity,
			nil,
		)
	}
	return fbAuthClient.client, nil
}

func (fbAuthClient *firebaseAuthenticationClient) InitializeClient(firebaseClient firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient) (bool, error) {
	app, err := firebaseClient.GetApp()
	if err != nil {
		return false, err
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			fbAuthClient.pathIdentity,
			err,
		)
	}
	fbAuthClient.firebaseServerlessClient = firebaseClient
	fbAuthClient.client = client

	return true, nil
}

func NewFirebaseAuthClient() (firebaseauthcoreclientinterfaces.FirebaseAuthenticationClient, error) {
	return &firebaseAuthenticationClient{
		pathIdentity: "FirebaseAuthClient",
	}, nil
}
