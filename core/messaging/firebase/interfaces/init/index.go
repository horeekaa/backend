package firebasemessagingcoreclientinterfaces

import (
	firebasemessagingcorewrapperinterfaces "github.com/horeekaa/backend/core/messaging/firebase/interfaces/wrappers"
	firebaseserverlesscoreclientinterfaces "github.com/horeekaa/backend/core/serverless/firebase/interfaces"
)

type FirebaseMessagingClient interface {
	InitializeClient(bServerlessClient firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient) (bool, error)
	GetServerlessClient() (firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient, error)
	GetMessagingClient() (firebasemessagingcorewrapperinterfaces.FirebaseMsgClient, error)
}
