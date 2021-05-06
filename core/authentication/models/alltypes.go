package authenticationcoremodels

import "firebase.google.com/go/v4/auth"

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

type AuthUserWrap struct {
	FirebaseUser *auth.UserRecord
}

type AuthTokenWrap struct {
	FirebaseToken *auth.Token
}
