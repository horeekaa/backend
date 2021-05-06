package firebaseauthdatasources

import (
	"context"

	auth "firebase.google.com/go/v4/auth"
	firebaseauthcoretypes "github.com/horeekaa/backend/core/authentication/firebase/types"
	firebaseauthcoreutilities "github.com/horeekaa/backend/core/authentication/firebase/utilities"
	authenticationcoreclientinterfaces "github.com/horeekaa/backend/core/authentication/interfaces"
	authenticationcoremodels "github.com/horeekaa/backend/core/authentication/models"

	firebaseauthcoreclientinterfaces "github.com/horeekaa/backend/core/authentication/firebase/interfaces"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
)

type firebaseAuthRepo struct {
	firebaseClient firebaseauthcoreclientinterfaces.FirebaseAuthenticationClient
}

func NewFirebaseAuthRepo(firebaseClient firebaseauthcoreclientinterfaces.FirebaseAuthenticationClient) (authenticationcoreclientinterfaces.AuthenticationRepo, error) {
	return &firebaseAuthRepo{
		firebaseClient,
	}, nil
}

func (fbAuthRepo *firebaseAuthRepo) VerifyAndDecodeToken(context context.Context, authToken string) (*authenticationcoremodels.AuthTokenWrap, error) {
	client, err := (*fbAuthRepo).firebaseClient.GetAuthClient()

	token, err := client.VerifyIDToken(context, authToken)
	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.DecodingAuthTokenFailed,
			"/authentication/verifyAndDecodeToken",
			err,
		)
	}
	return &authenticationcoremodels.AuthTokenWrap{
		FirebaseToken: token,
	}, nil
}

func (fbAuthRepo *firebaseAuthRepo) GetAuthUserDataByEmail(context context.Context, email string) (*authenticationcoremodels.AuthUserWrap, error) {
	client, err := (*fbAuthRepo).firebaseClient.GetAuthClient()

	user, err := client.GetUserByEmail(context, email)
	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.GetAuthDataFailed,
			"/authentication/getAuthUserDataByEmail",
			err,
		)
	}
	return &authenticationcoremodels.AuthUserWrap{
		FirebaseUser: user,
	}, nil
}

func (fbAuthRepo *firebaseAuthRepo) GetAuthUserDataById(context context.Context, uid string) (*authenticationcoremodels.AuthUserWrap, error) {
	client, err := (*fbAuthRepo).firebaseClient.GetAuthClient()

	user, err := client.GetUser(context, uid)
	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.GetAuthDataFailed,
			"/authentication/getAuthUserDataById",
			err,
		)
	}
	return &authenticationcoremodels.AuthUserWrap{
		FirebaseUser: user,
	}, nil
}

func (fbAuthRepo *firebaseAuthRepo) SetRoleInAuthUserData(context context.Context, uid string, accountType string, dbID string) (bool, error) {
	claims := map[string]interface{}{
		firebaseauthcoretypes.FirebaseCustomClaimsAccountTypeKey: accountType,
		firebaseauthcoretypes.FirebaseCustomClaimsAccountIDKey:   dbID,
	}
	client, _ := (*fbAuthRepo).firebaseClient.GetAuthClient()

	if err := client.SetCustomUserClaims(context, uid, claims); err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.SetAuthDataFailed,
			"/authentication/setRoleInAuthUserData",
			err,
		)
	}

	return true, nil
}

func (fbAuthRepo *firebaseAuthRepo) UpdateAuthUserData(context context.Context, user *authenticationcoremodels.UpdateAuthUserData) (*authenticationcoremodels.AuthUserWrap, error) {
	params := (&auth.UserToUpdate{})
	if &user.Email != nil {
		params = params.Email((*user).Email)
	}
	if &user.EmailVerified != nil {
		params = params.EmailVerified((*user).EmailVerified)
	}
	if &user.PhoneNumber != nil {
		params = params.PhoneNumber((*user).PhoneNumber)
	}
	if &user.Password != nil {
		params = params.Password((*user).Password)
	}
	if &user.DisplayName != nil {
		params = params.DisplayName((*user).DisplayName)
	}
	if &user.PhotoURL != nil {
		params = params.PhotoURL((*user).PhotoURL)
	}
	if &user.Disabled != nil {
		params = params.Disabled((*user).Disabled)
	}
	client, err := (*fbAuthRepo).firebaseClient.GetAuthClient()

	updatedUser, err := client.UpdateUser(context, (*user).UID, params)
	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpdateObjectFailed,
			"/authentication/updateAuthUserData",
			err,
		)
	}
	return &authenticationcoremodels.AuthUserWrap{
		FirebaseUser: updatedUser,
	}, nil
}

func (fbAuthRepo *firebaseAuthRepo) GenerateEmailVerificationLink(context context.Context, email string) (string, error) {
	client, err := (*fbAuthRepo).firebaseClient.GetAuthClient()

	link, err := client.EmailVerificationLinkWithSettings(
		context,
		email,
		firebaseauthcoreutilities.GetFirebaseActionCodeSettings()["data"].(*auth.ActionCodeSettings),
	)
	if err != nil {
		return "", horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			"firebase/generateEmailVerificationLink",
			err,
		)
	}
	return link, nil
}

func (fbAuthRepo *firebaseAuthRepo) GeneratePasswordResetLink(context context.Context, email string) (string, error) {
	client, err := (*fbAuthRepo).firebaseClient.GetAuthClient()

	link, err := client.PasswordResetLinkWithSettings(
		context,
		email,
		firebaseauthcoreutilities.GetFirebaseActionCodeSettings()["data"].(*auth.ActionCodeSettings),
	)
	if err != nil {
		return "", horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			"firebase/passwordResetLinkWithSettings",
			err,
		)
	}
	return link, nil
}
