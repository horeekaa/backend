package firebaseserverlesscorewrapperinterfaces

import (
	"context"

	firebaseauthcorewrapperinterfaces "github.com/horeekaa/backend/core/authentication/firebase/interfaces/wrappers"
)

type FirebaseApp interface {
	Auth(ctx context.Context) (firebaseauthcorewrapperinterfaces.FirebaseAuthClient, error)
}
