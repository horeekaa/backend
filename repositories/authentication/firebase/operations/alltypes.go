package firebaseauthoperations

import (
	auth "firebase.google.com/go/v4/auth"
)

type FirebaseAuthToken struct {
	Token *auth.Token
}

type FirebaseUserRecord struct {
	User *auth.UserRecord
}
