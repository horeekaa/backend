package firebaseauthcorewrapperinterfaces

import (
	"context"

	auth "firebase.google.com/go/v4/auth"
)

type FirebaseAuthClient interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
	GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error)
	GetUser(ctx context.Context, uid string) (*auth.UserRecord, error)
	SetCustomUserClaims(ctx context.Context, uid string, customClaims map[string]interface{}) error
	UpdateUser(ctx context.Context, uid string, user *auth.UserToUpdate) (ur *auth.UserRecord, err error)
	EmailVerificationLinkWithSettings(ctx context.Context, email string, settings *auth.ActionCodeSettings) (string, error)
	PasswordResetLinkWithSettings(ctx context.Context, email string, settings *auth.ActionCodeSettings) (string, error)
}
