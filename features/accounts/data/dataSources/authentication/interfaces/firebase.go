package firebaseauthdatasourceinterfaces

import (
	"context"

	auth "firebase.google.com/go/v4/auth"
	authenticationcoremodels "github.com/horeekaa/backend/core/authentication/models"
)

type FirebaseAuthRepo interface {
	VerifyAndDecodeToken(context context.Context, authToken string) (*auth.Token, error)
	GetAuthUserDataByEmail(context context.Context, email string) (*auth.UserRecord, error)
	GetAuthUserDataById(context context.Context, uid string) (*auth.UserRecord, error)
	SetRoleInAuthUserData(context context.Context, uid string, accountType string, dbID string) (bool, error)
	UpdateAuthUserData(context context.Context, user *authenticationcoremodels.UpdateAuthUserData) (*auth.UserRecord, error)
	GenerateEmailVerificationLink(context context.Context, email string) (string, error)
	GeneratePasswordResetLink(context context.Context, email string) (string, error)
}
