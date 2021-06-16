package firebaseserverlesscorewrappers

import (
	"context"

	firebase "firebase.google.com/go/v4"
	firebaseauthcorewrapperinterfaces "github.com/horeekaa/backend/core/authentication/firebase/interfaces/wrappers"
	firebaseserverlesscorewrapperinterfaces "github.com/horeekaa/backend/core/serverless/firebase/interfaces/wrappers"
)

type firebaseApp struct {
	*firebase.App
}

func (fbApp *firebaseApp) Auth(ctx context.Context) (firebaseauthcorewrapperinterfaces.FirebaseAuthClient, error) {
	return fbApp.App.Auth(ctx)
}

func NewFirebaseApp(wrappedFirebaseApp *firebase.App) (firebaseserverlesscorewrapperinterfaces.FirebaseApp, error) {
	return &firebaseApp{
		wrappedFirebaseApp,
	}, nil
}
