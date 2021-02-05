package firebaseauthinterfaces

import (
	auth "firebase.google.com/go/v4/auth"
	firebaseauthmodels "github.com/horeekaa/backend/repositories/authentication/firebase/models"
)

type FirebaseAuthentication interface {
	VerifyAndDecodeToken(authToken string) (*auth.Token, error)
	GetAuthUserDataByEmail(email string) (*auth.UserRecord, error)
	GetAuthUserDataById(uid string) (*auth.UserRecord, error)
	SetRoleInAuthUserData(user *auth.UserRecord, accountRole string, dbId string) (bool, error)
	UpdateAuthUserData(user *firebaseauthmodels.UpdateAuthUserData) (*auth.UserRecord, error)
	GenerateEmailVerificationLink(email string) (string, error)
	GeneratePasswordResetLink(email string) (string, error)
}
