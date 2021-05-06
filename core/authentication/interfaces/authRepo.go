package authenticationcoreclientinterfaces

import (
	"context"

	authenticationcoremodels "github.com/horeekaa/backend/core/authentication/models"
)

type AuthenticationRepo interface {
	VerifyAndDecodeToken(context context.Context, authToken string) (*authenticationcoremodels.AuthTokenWrap, error)
	GetAuthUserDataByEmail(context context.Context, email string) (*authenticationcoremodels.AuthUserWrap, error)
	GetAuthUserDataById(context context.Context, uid string) (*authenticationcoremodels.AuthUserWrap, error)
	SetRoleInAuthUserData(context context.Context, uid string, accountType string, dbID string) (bool, error)
	UpdateAuthUserData(context context.Context, user *authenticationcoremodels.UpdateAuthUserData) (*authenticationcoremodels.AuthUserWrap, error)
	GenerateEmailVerificationLink(context context.Context, email string) (string, error)
	GeneratePasswordResetLink(context context.Context, email string) (string, error)
}
