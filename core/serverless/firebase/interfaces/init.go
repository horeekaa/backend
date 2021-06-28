package firebaseserverlesscoreclientinterfaces

import (
	firebaseserverlesscorewrapperinterfaces "github.com/horeekaa/backend/core/serverless/firebase/interfaces/wrappers"
)

type FirebaseServerlessClient interface {
	Connect() (bool, error)
	GetApp() (firebaseserverlesscorewrapperinterfaces.FirebaseApp, error)
}
