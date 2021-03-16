package firebaseauthcoreinterfaces

import (
	firebaseauthcoremodels "github.com/horeekaa/backend/core/authentication/firebase/models"
	firebaseauthcoreoperations "github.com/horeekaa/backend/core/authentication/firebase/operations"
)

type FirebaseAuthentication interface {
	VerifyAndDecodeToken(authToken string) (*firebaseauthcoreoperations.FirebaseAuthToken, error)
	GetAuthUserDataByEmail(email string) (*firebaseauthcoreoperations.FirebaseUserRecord, error)
	GetAuthUserDataById(uid string) (*firebaseauthcoreoperations.FirebaseUserRecord, error)
	SetRoleInAuthUserData(uid string, accountType string, dbID string) (bool, error)
	UpdateAuthUserData(user *firebaseauthcoremodels.UpdateAuthUserData) (*firebaseauthcoreoperations.FirebaseUserRecord, error)
	GenerateEmailVerificationLink(email string) (string, error)
	GeneratePasswordResetLink(email string) (string, error)
}
