package firebaseserverlesscorewrapperinterfaces

import (
	"context"

	firebaseauthcorewrapperinterfaces "github.com/horeekaa/backend/core/authentication/firebase/interfaces/wrappers"
	firebasemessagingcorewrapperinterfaces "github.com/horeekaa/backend/core/messaging/firebase/interfaces/wrappers"
)

type FirebaseApp interface {
	Auth(ctx context.Context) (firebaseauthcorewrapperinterfaces.FirebaseAuthClient, error)
	Messaging(ctx context.Context) (firebasemessagingcorewrapperinterfaces.FirebaseMsgClient, error)
}
