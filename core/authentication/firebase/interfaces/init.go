package firebaseauthcoreinterfaces

import (
	firebaseauthcoremodels "github.com/horeekaa/backend/core/authentication/firebase/models"
	firebaseauthcoretypes "github.com/horeekaa/backend/core/authentication/firebase/types"
)

type FirebaseAuthentication interface {
	VerifyAndDecodeToken(authToken string) (*firebaseauthcoretypes.FirebaseAuthToken, error)
	GetAuthUserDataByEmail(email string) (*firebaseauthcoretypes.FirebaseUserRecord, error)
	GetAuthUserDataById(uid string) (*firebaseauthcoretypes.FirebaseUserRecord, error)
	SetRoleInAuthUserData(uid string, accountType string, dbID string) (bool, error)
	UpdateAuthUserData(user *firebaseauthcoremodels.UpdateAuthUserData) (*firebaseauthcoretypes.FirebaseUserRecord, error)
	GenerateEmailVerificationLink(email string) (string, error)
	GeneratePasswordResetLink(email string) (string, error)
}
