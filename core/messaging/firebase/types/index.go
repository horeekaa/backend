package firebasemessagingcoretypes

import "firebase.google.com/go/v4/messaging"

type SentMessage struct {
	*messaging.Message
}

type SentMulticastMessage struct {
	*messaging.MulticastMessage
}

type BatchMessageResponse struct {
	*messaging.BatchResponse
}
