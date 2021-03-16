package firebaseauthcoreclients

import (
	"context"

	firebase "firebase.google.com/go/v4"
	auth "firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	horeekaaexception "github.com/horeekaa/backend/_errors/repoExceptions"
	horeekaaexceptionenums "github.com/horeekaa/backend/_errors/repoExceptions/_enums"
	coreconfigs "github.com/horeekaa/backend/core/_commons/configs"
	firebaseauthcoreinterfaces "github.com/horeekaa/backend/core/authentication/firebase/interfaces"
	firebaseauthcoremodels "github.com/horeekaa/backend/core/authentication/firebase/models"
	firebaseauthcoreoperations "github.com/horeekaa/backend/core/authentication/firebase/operations"
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
		return nil, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.UpstreamException,
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

func (fbAuth *firebaseAuthentication) VerifyAndDecodeToken(authToken string) (*firebaseauthcoreoperations.FirebaseAuthToken, error) {
	token, err := (*fbAuth).Client.VerifyIDToken(*fbAuth.Context, authToken)
	if err != nil {
		return nil, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.DecodingAuthTokenFailed,
			"/authentication/verifyAndDecodeToken",
			err,
		)
	}
	return &firebaseauthcoreoperations.FirebaseAuthToken{
		Token: token,
	}, nil
}

func (fbAuth *firebaseAuthentication) GetAuthUserDataByEmail(email string) (*firebaseauthcoreoperations.FirebaseUserRecord, error) {
	user, err := (*fbAuth).Client.GetUserByEmail(*fbAuth.Context, email)
	if err != nil {
		return nil, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.GetAuthDataFailed,
			"/authentication/getAuthUserDataByEmail",
			err,
		)
	}
	return &firebaseauthcoreoperations.FirebaseUserRecord{
		RepoUser: user,
	}, nil
}

func (fbAuth *firebaseAuthentication) GetAuthUserDataById(uid string) (*firebaseauthcoreoperations.FirebaseUserRecord, error) {
	user, err := (*fbAuth).Client.GetUser(*fbAuth.Context, uid)
	if err != nil {
		return nil, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.GetAuthDataFailed,
			"/authentication/getAuthUserDataById",
			err,
		)
	}
	return &firebaseauthcoreoperations.FirebaseUserRecord{
		RepoUser: user,
	}, nil
}

func (fbAuth *firebaseAuthentication) SetRoleInAuthUserData(uid string, accountType string, dbID string) (bool, error) {
	claims := map[string]interface{}{
		firebaseauthcoreoperations.FirebaseCustomClaimsAccountTypeKey: accountType,
		firebaseauthcoreoperations.FirebaseCustomClaimsAccountIDKey:   dbID,
	}
	if err := (*fbAuth).Client.SetCustomUserClaims(*fbAuth.Context, uid, claims); err != nil {
		return false, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.SetAuthDataFailed,
			"/authentication/setRoleInAuthUserData",
			err,
		)
	}

	return true, nil
}

func (fbAuth *firebaseAuthentication) UpdateAuthUserData(user *firebaseauthcoremodels.UpdateAuthUserData) (*firebaseauthcoreoperations.FirebaseUserRecord, error) {
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
		return nil, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.UpdateObjectFailed,
			"/authentication/updateAuthUserData",
			err,
		)
	}
	return &firebaseauthcoreoperations.FirebaseUserRecord{
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
		return "", horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.UpstreamException,
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
		return "", horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.UpstreamException,
			"firebase/passwordResetLinkWithSettings",
			err,
		)
	}
	return link, nil
}
