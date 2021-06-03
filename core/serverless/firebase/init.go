package firebaseserverlesscoreclients

import (
	"context"

	firebase "firebase.google.com/go/v4"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	firebaseserverlesscoreclientinterfaces "github.com/horeekaa/backend/core/serverless/firebase/interfaces"
)

type firebaseServerlessClient struct {
	app *firebase.App
}

func (svlessClient *firebaseServerlessClient) Connect() (bool, error) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			"/newFirebaseAuthentication",
			err,
		)
	}

	svlessClient.app = app
	return true, nil
}

func (svlessClient *firebaseServerlessClient) GetApp() (*firebase.App, error) {
	if svlessClient.app == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClientInitializationFailed,
			"/newMongoClient",
			nil,
		)
	}
	return svlessClient.app, nil
}

func NewFirebaseServerlessClient() (firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient, error) {
	return &firebaseServerlessClient{}, nil
}
