package firebaseauthinterfaces

import (
	firebaseauthmodels "github.com/horeekaa/backend/repositories/authentication/firebase/models"
)

type FirebaseAuthentication interface {
	VerifyAndDecodeToken(authToken string) (*firebaseauthmodels.FirebaseAuthToken, error)
	GetAuthUserDataByEmail(email string) (*firebaseauthmodels.FirebaseUserRecord, error)
	GetAuthUserDataById(uid string) (*firebaseauthmodels.FirebaseUserRecord, error)
	SetRoleInAuthUserData(user *firebaseauthmodels.FirebaseUserRecord, accountRole string, dbId string) (bool, error)
	UpdateAuthUserData(user *firebaseauthmodels.UpdateAuthUserData) (*firebaseauthmodels.FirebaseUserRecord, error)
	GenerateEmailVerificationLink(email string) (string, error)
	GeneratePasswordResetLink(email string) (string, error)
}
