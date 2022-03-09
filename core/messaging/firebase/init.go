package firebasemessagingcoreclients

import (
	"context"

	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	firebasemessagingcoreclientinterfaces "github.com/horeekaa/backend/core/messaging/firebase/interfaces/init"
	firebasemessagingcorewrapperinterfaces "github.com/horeekaa/backend/core/messaging/firebase/interfaces/wrappers"
	firebaseserverlesscoreclientinterfaces "github.com/horeekaa/backend/core/serverless/firebase/interfaces"
)

type firebaseMessagingClient struct {
	firebaseServerlessClient firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient
	client                   firebasemessagingcorewrapperinterfaces.FirebaseMsgClient
	pathIdentity string
}

func (fbMsgClient *firebaseMessagingClient) GetServerlessClient() (firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient, error) {
	if fbMsgClient.firebaseServerlessClient == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClientInitializationFailed,
			fbMsgClient.pathIdentity,
			nil,
		)
	}
	return fbMsgClient.firebaseServerlessClient, nil
}

func (fbMsgClient *firebaseMessagingClient) GetMessagingClient() (firebasemessagingcorewrapperinterfaces.FirebaseMsgClient, error) {
	if fbMsgClient.client == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.ClientInitializationFailed,
			fbMsgClient.pathIdentity,
			nil,
		)
	}
	return fbMsgClient.client, nil
}

func (fbMsgClient *firebaseMessagingClient) InitializeClient(firebaseClient firebaseserverlesscoreclientinterfaces.FirebaseServerlessClient) (bool, error) {
	app, err := firebaseClient.GetApp()
	if err != nil {
		return false, err
	}
	client, err := app.Messaging(context.Background())
	if err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			fbMsgClient.pathIdentity,
			err,
		)
	}

	fbMsgClient.firebaseServerlessClient = firebaseClient
	fbMsgClient.client = client

	return true, nil
}

func NewFirebaseMsgClient() (firebasemessagingcoreclientinterfaces.FirebaseMessagingClient, error) {
	return &firebaseMessagingClient{
		pathIdentity: "FirebaseMessagingClient",
	}, nil
}
