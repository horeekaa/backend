package firebaseauthoperations

import (
	auth "firebase.google.com/go/v4/auth"
)

const (
	FirebaseCustomClaimsAccountIDKey   = "AccountId"
	FirebaseCustomClaimsAccountTypeKey = "Type"
)

type FirebaseAuthToken struct {
	Token *auth.Token
}

type FirebaseUserRecord struct {
	FirebaseUser *auth.UserRecord
}
