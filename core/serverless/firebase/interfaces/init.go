package firebaseserverlesscoreclientinterfaces

import (
	firebase "firebase.google.com/go/v4"
)

type FirebaseServerlessClient interface {
	Connect() (bool, error)
	GetApp() (*firebase.App, error)
}
