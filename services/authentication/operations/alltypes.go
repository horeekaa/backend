package authserviceoperations

import (
	firebaseauthoperations "github.com/horeekaa/backend/repositories/authentication/firebase/operations"
)

const (
	VerifyEmailAuthURLAction   string = "VERIFY_EMAIL_AUTH_URL_ACTION"
	ResetPasswordAuthURLAction string = "RESET_PASSWORD_AUTH_URL_ACTION"
)

type AuthenticationServiceToken struct {
	ServiceToken *firebaseauthoperations.FirebaseAuthToken
}

type AuthenticationServiceUserRecord struct {
	ServiceUser *firebaseauthoperations.FirebaseUserRecord
}
