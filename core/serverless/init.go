package serverlesscoreclients

import (
	firebaseserverlesscoreclientinterfaces "github.com/horeekaa/backend/core/serverless/firebase/interfaces"
	serverlesscoreclientinterfaces "github.com/horeekaa/backend/core/serverless/interfaces"
)

type serverlessClient struct {
	firebaseApp firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient
}

func (svlessClient *serverlessClient) SetFirebaseApp(firebaseApp firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient) bool {
	svlessClient.firebaseApp = firebaseApp
	return true
}

func (svlessClient *serverlessClient) GetFirebaseApp() firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient {
	return svlessClient.firebaseApp
}

func NewServerlessClient() (serverlesscoreclientinterfaces.ServerlessClient, error) {
	return &serverlessClient{}, nil
}
