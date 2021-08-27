package firebasemessagingcoreoperationinterfaces

import (
	"context"

	firebasemessagingcoretypes "github.com/horeekaa/backend/core/messaging/firebase/types"
)

type FirebaseMessagingBasicOperation interface {
	SendMessage(ctx context.Context, message *firebasemessagingcoretypes.SentMessage) (string, error)
	SendMulticastMessage(ctx context.Context, message *firebasemessagingcoretypes.SentMulticastMessage) (*firebasemessagingcoretypes.BatchMessageResponse, error)
}
