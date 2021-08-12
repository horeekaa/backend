package messagingcoredependencies

import (
	"github.com/golobby/container/v2"
	firebasemessagingcoreclients "github.com/horeekaa/backend/core/messaging/firebase"
	firebasemessagingcoreclientinterfaces "github.com/horeekaa/backend/core/messaging/firebase/interfaces/init"
	firebasemessagingcoreoperationinterfaces "github.com/horeekaa/backend/core/messaging/firebase/interfaces/operations"
	firebasemessagingcoreoperations "github.com/horeekaa/backend/core/messaging/firebase/operations"
	firebaseserverlesscoreclientinterfaces "github.com/horeekaa/backend/core/serverless/firebase/interfaces"
)

type FirebaseMessagingDependency struct{}

func (_ FirebaseMessagingDependency) Bind() {
	container.Singleton(
		func(fbSvless firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient) firebasemessagingcoreclientinterfaces.FirebaseMessagingClient {
			fbMsgClient, _ := firebasemessagingcoreclients.NewFirebaseMsgClient()
			fbMsgClient.InitializeClient(fbSvless)
			return fbMsgClient
		},
	)

	container.Singleton(
		func(
			fbMsgClient firebasemessagingcoreclientinterfaces.FirebaseMessagingClient,
		) firebasemessagingcoreoperationinterfaces.FirebaseMessagingBasicOperation {
			fbMsgBasicOps, _ := firebasemessagingcoreoperations.NewFirebaseMessagingBasicOperation(
				fbMsgClient,
			)
			return fbMsgBasicOps
		},
	)
}
