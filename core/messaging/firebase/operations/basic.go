package firebasemessagingcoreoperations

import (
	"context"

	"firebase.google.com/go/v4/messaging"
	firebasemessagingcoreclientinterfaces "github.com/horeekaa/backend/core/messaging/firebase/interfaces/init"
	firebasemessagingcoreoperationinterfaces "github.com/horeekaa/backend/core/messaging/firebase/interfaces/operations"
	firebasemessagingcoretypes "github.com/horeekaa/backend/core/messaging/firebase/types"
)

type firebaseMessagingBasicOperation struct {
	firebaseMsgClient firebasemessagingcoreclientinterfaces.FirebaseMessagingClient
	pathIdentity      string
}

func NewFirebaseMessagingBasicOperation(firebaseMsgClient firebasemessagingcoreclientinterfaces.FirebaseMessagingClient) (firebasemessagingcoreoperationinterfaces.FirebaseMessagingBasicOperation, error) {
	return &firebaseMessagingBasicOperation{
		firebaseMsgClient,
		"FirebaseMessagingBasicOperation",
	}, nil
}

func (fbMsgBasicOps *firebaseMessagingBasicOperation) SendMessage(ctx context.Context, message *firebasemessagingcoretypes.SentMessage) (string, error) {
	client, _ := fbMsgBasicOps.firebaseMsgClient.GetMessagingClient()

	nativeMessage := messaging.Message(*message)
	res, _ := client.Send(ctx, &nativeMessage)

	return res, nil
}

func (fbMsgBasicOps *firebaseMessagingBasicOperation) SendMulticastMessage(ctx context.Context, message *firebasemessagingcoretypes.SentMulticastMessage) (*firebasemessagingcoretypes.BatchMessageResponse, error) {
	client, _ := fbMsgBasicOps.firebaseMsgClient.GetMessagingClient()

	nativeMessage := messaging.MulticastMessage(*message)
	_, _ = client.SendMulticast(ctx, &nativeMessage)

	batchResponse := firebasemessagingcoretypes.BatchMessageResponse{}
	return &batchResponse, nil
}
