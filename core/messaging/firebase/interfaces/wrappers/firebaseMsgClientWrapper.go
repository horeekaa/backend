package firebasemessagingcorewrapperinterfaces

import (
	"context"

	messaging "firebase.google.com/go/v4/messaging"
)

type FirebaseMsgClient interface {
	Send(ctx context.Context, message *messaging.Message) (string, error)
	SendAll(ctx context.Context, messages []*messaging.Message) (*messaging.BatchResponse, error)
	SendMulticast(ctx context.Context, message *messaging.MulticastMessage) (*messaging.BatchResponse, error)
	SubscribeToTopic(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error)
	UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error)
}
