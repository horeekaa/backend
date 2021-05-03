package serverlesscoredependencies

import (
	"github.com/golobby/container/v2"
	serverlesscoreclients "github.com/horeekaa/backend/core/serverless"
	firebaseserverlesscoreclients "github.com/horeekaa/backend/core/serverless/firebase"
	firebaseserverlesscoreclientinterfaces "github.com/horeekaa/backend/core/serverless/firebase/interfaces"
	serverlesscoreclientinterfaces "github.com/horeekaa/backend/core/serverless/interfaces"
)

type FirebaseServerlessDependency struct{}

func (firebaseDependency *FirebaseServerlessDependency) Bind() {
	container.Singleton(
		func() firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient {
			fbServerlessClient, _ := firebaseserverlesscoreclients.NewFirebaseServerlessClient()
			fbServerlessClient.Connect()
			return fbServerlessClient
		},
	)

	container.Singleton(
		func(firebaseServerlessClient firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient) serverlesscoreclientinterfaces.ServerlessClient {
			serverlessClient, _ := serverlesscoreclients.NewServerlessClient()
			serverlessClient.SetFirebaseApp(firebaseServerlessClient)

			return serverlessClient
		},
	)
}
