package firebaseserverlesscoreclients

import (
	"context"
	"errors"

	firebase "firebase.google.com/go/v4"
	coreconfigs "github.com/horeekaa/backend/core/_commons/configs"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	firebaseserverlesscoreclientinterfaces "github.com/horeekaa/backend/core/serverless/firebase/interfaces"
	"google.golang.org/api/option"
)

type firebaseServerlessClient struct {
	app *firebase.App
}

func (svlessClient *firebaseServerlessClient) Connect() (bool, error) {
	opt := option.WithCredentialsFile(coreconfigs.GetEnvVariable(coreconfigs.FirebaseServiceAccountPath))
	config := &firebase.Config{ProjectID: coreconfigs.GetEnvVariable(coreconfigs.FirebaseConfigProjectID)}
	app, err := firebase.NewApp(context.Background(), config, opt)
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
			errors.New(horeekaacoreexceptionenums.ClientInitializationFailed),
		)
	}
	return svlessClient.app, nil
}

func NewFirebaseServerlessClient() (firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient, error) {
	return &firebaseServerlessClient{}, nil
}
