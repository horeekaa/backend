package firebaseserverlesscoreclients

import (
	"context"

	firebase "firebase.google.com/go/v4"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	firebaseserverlesscoreclientinterfaces "github.com/horeekaa/backend/core/serverless/firebase/interfaces"
	firebaseserverlesscorewrapperinterfaces "github.com/horeekaa/backend/core/serverless/firebase/interfaces/wrappers"
	firebaseserverlesscorewrappers "github.com/horeekaa/backend/core/serverless/firebase/wrappers"
)

type firebaseServerlessClient struct {
	app firebaseserverlesscorewrapperinterfaces.FirebaseApp
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

	wrappedApp, _ := firebaseserverlesscorewrappers.NewFirebaseApp(app)
	svlessClient.app = wrappedApp
	return true, nil
}

func (svlessClient *firebaseServerlessClient) GetApp() (firebaseserverlesscorewrapperinterfaces.FirebaseApp, error) {
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
