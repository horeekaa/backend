package firebaseauthinterfaces

import (
	firebaseauthmodels "github.com/horeekaa/backend/repositories/authentication/firebase/models"
	firebaseauthoperations "github.com/horeekaa/backend/repositories/authentication/firebase/operations"
)

type FirebaseAuthentication interface {
	VerifyAndDecodeToken(authToken string) (*firebaseauthoperations.FirebaseAuthToken, error)
	GetAuthUserDataByEmail(email string) (*firebaseauthoperations.FirebaseUserRecord, error)
	GetAuthUserDataById(uid string) (*firebaseauthoperations.FirebaseUserRecord, error)
	SetRoleInAuthUserData(uid string, accountType string, dbID string) (bool, error)
	UpdateAuthUserData(user *firebaseauthmodels.UpdateAuthUserData) (*firebaseauthoperations.FirebaseUserRecord, error)
	GenerateEmailVerificationLink(email string) (string, error)
	GeneratePasswordResetLink(email string) (string, error)
}
