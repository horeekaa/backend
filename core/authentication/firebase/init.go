package firebaseauthcoreclients

import (
	"context"

	firebase "firebase.google.com/go/v4"
	auth "firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	coreconfigs "github.com/horeekaa/backend/core/_commons/configs"
	horeekaacoreexception "github.com/horeekaa/backend/core/_errors/repoExceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/_errors/repoExceptions/_enums"
	firebaseauthcoreinterfaces "github.com/horeekaa/backend/core/authentication/firebase/interfaces"
	firebaseauthcoremodels "github.com/horeekaa/backend/core/authentication/firebase/models"
	firebaseauthcoretypes "github.com/horeekaa/backend/core/authentication/firebase/types"
	firebaseauthcoreutilities "github.com/horeekaa/backend/core/authentication/firebase/utilities"
)

type firebaseAuthentication struct {
	App     *firebase.App
	Client  *auth.Client
	Context *context.Context
}

func NewFirebaseAuthentication(context *context.Context) (firebaseauthcoreinterfaces.FirebaseAuthentication, error) {
	opt := option.WithCredentialsFile(coreconfigs.GetEnvVariable(coreconfigs.FirebaseServiceAccountPath))
	config := &firebase.Config{ProjectID: coreconfigs.GetEnvVariable(coreconfigs.FirebaseConfigProjectID)}
	app, err := firebase.NewApp(*context, config, opt)
	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpstreamException,
			"/newFirebaseAuthentication",
			err,
		)
	}

	client, err := app.Auth(*context)
	return &firebaseAuthentication{
		App:     app,
		Client:  client,
		Context: context,
	}, nil
}

func (fbAuth *firebaseAuthentication) VerifyAndDecodeToken(authToken string) (*firebaseauthcoretypes.FirebaseAuthToken, error) {
	token, err := (*fbAuth).Client.VerifyIDToken(*fbAuth.Context, authToken)
	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.DecodingAuthTokenFailed,
			"/authentication/verifyAndDecodeToken",
			err,
		)
	}
	return &firebaseauthcoretypes.FirebaseAuthToken{
		Token: token,
	}, nil
}

func (fbAuth *firebaseAuthentication) GetAuthUserDataByEmail(email string) (*firebaseauthcoretypes.FirebaseUserRecord, error) {
	user, err := (*fbAuth).Client.GetUserByEmail(*fbAuth.Context, email)
	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.GetAuthDataFailed,
			"/authentication/getAuthUserDataByEmail",
			err,
		)
	}
	return &firebaseauthcoretypes.FirebaseUserRecord{
		RepoUser: user,
	}, nil
}

func (fbAuth *firebaseAuthentication) GetAuthUserDataById(uid string) (*firebaseauthcoretypes.FirebaseUserRecord, error) {
	user, err := (*fbAuth).Client.GetUser(*fbAuth.Context, uid)
	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.GetAuthDataFailed,
			"/authentication/getAuthUserDataById",
			err,
		)
	}
	return &firebaseauthcoretypes.FirebaseUserRecord{
		RepoUser: user,
	}, nil
}

func (fbAuth *firebaseAuthentication) SetRoleInAuthUserData(uid string, accountType string, dbID string) (bool, error) {
	claims := map[string]interface{}{
		firebaseauthcoretypes.FirebaseCustomClaimsAccountTypeKey: accountType,
		firebaseauthcoretypes.FirebaseCustomClaimsAccountIDKey:   dbID,
	}
	if err := (*fbAuth).Client.SetCustomUserClaims(*fbAuth.Context, uid, claims); err != nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.SetAuthDataFailed,
			"/authentication/setRoleInAuthUserData",
			err,
		)
	}

	return true, nil
}

func (fbAuth *firebaseAuthentication) UpdateAuthUserData(user *firebaseauthcoremodels.UpdateAuthUserData) (*firebaseauthcoretypes.FirebaseUserRecord, error) {
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

	updatedUser, err := (*fbAuth).Client.UpdateUser(*fbAuth.Context, (*user).UID, params)
	if err != nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.UpdateObjectFailed,
			"/authentication/updateAuthUserData",
			err,
		)
	}
	return &firebaseauthcoretypes.FirebaseUserRecord{
		RepoUser: updatedUser,
	}, nil
}

func (fbAuth *firebaseAuthentication) GenerateEmailVerificationLink(email string) (string, error) {
	link, err := (*fbAuth).Client.EmailVerificationLinkWithSettings(
		*fbAuth.Context,
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

func (fbAuth *firebaseAuthentication) GeneratePasswordResetLink(email string) (string, error) {
	link, err := (*fbAuth).Client.PasswordResetLinkWithSettings(
		*fbAuth.Context,
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
