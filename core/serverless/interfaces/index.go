package serverlesscoreclientinterfaces

import (
	firebaseserverlesscoreclientinterfaces "github.com/horeekaa/backend/core/serverless/firebase/interfaces"
)

type ServerlessClient interface {
	SetFirebaseApp(firebaseApp firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient) bool
	GetFirebaseApp() firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient
}
