package firebasemessagingcoreoperations

import (
	"context"

	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	firebasemessagingcoreclientinterfaces "github.com/horeekaa/backend/core/messaging/firebase/interfaces/init"
	firebasemessagingcoreoperationinterfaces "github.com/horeekaa/backend/core/messaging/firebase/interfaces/operations"
	firebasemessagingcoretypes "github.com/horeekaa/backend/core/messaging/firebase/types"
)

type firebaseMessagingBasicOperation struct {
	firebaseMsgClient firebasemessagingcoreclientinterfaces.FirebaseMessagingClient
}

func NewFirebaseMessagingBasicOperation(firebaseMsgClient firebasemessagingcoreclientinterfaces.FirebaseMessagingClient) (firebasemessagingcoreoperationinterfaces.FirebaseMessagingBasicOperation, error) {
	return &firebaseMessagingBasicOperation{
		firebaseMsgClient,
	}, nil
}

func (fbMsgBasicOps *firebaseMessagingBasicOperation) SendMessage(ctx context.Context, message *firebasemessagingcoretypes.SentMessage) (string, error) {
	client, _ := fbMsgBasicOps.firebaseMsgClient.GetMessagingClient()

	res, err := client.Send(ctx, message.Message)
	if err != nil {
		return "", horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.SendNotifMessageFailed,
			"/fbMsgBasicOperation/SendMessage",
			err,
		)
	}

	return res, nil
}

func (fbMsgBasicOps *firebaseMessagingBasicOperation) SendMulticastMessage(ctx context.Context, message *firebasemessagingcoretypes.SentMulticastMessage) (*firebasemessagingcoretypes.BatchMessageResponse, error) {
	client, _ := fbMsgBasicOps.firebaseMsgClient.GetMessagingClient()

	res, err := client.SendMulticast(ctx, message.MulticastMessage)
	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.SendNotifMessageFailed,
			"/fbMsgBasicOperation/SendMulticastMessage",
			err,
		)
	}

	return &firebasemessagingcoretypes.BatchMessageResponse{
		BatchResponse: res,
	}, nil
}
