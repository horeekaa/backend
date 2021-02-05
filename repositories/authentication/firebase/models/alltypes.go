package firebaseauthmodels

import (
	auth "firebase.google.com/go/v4/auth"
)

type UpdateAuthUserData struct {
	UID           string
	Email         string
	EmailVerified bool
	PhoneNumber   string
	Password      string
	DisplayName   string
	PhotoURL      string
	Disabled      bool
}

type FirebaseAuthToken struct {
	Token *auth.Token
}

type FirebaseUserRecord struct {
	User *auth.UserRecord
}
