package firebaseauthentication

import (
	"context"

	firebase "firebase.google.com/go/v4"
	auth "firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	configs "github.com/horeekaa/backend/_commons/configs"
	horeekaaexception "github.com/horeekaa/backend/_errors/repoExceptions"
	horeekaaexceptionenums "github.com/horeekaa/backend/_errors/repoExceptions/_enums"
	firebaseauthenticationinterface "github.com/horeekaa/backend/repositories/authentication/firebase/interfaces"
	firebaseauthenticationmodel "github.com/horeekaa/backend/repositories/authentication/firebase/models"
)

type firebaseAuthentication struct {
	App     *firebase.App
	Client  *auth.Client
	Context *context.Context
}

func NewFirebaseAuthentication(context *context.Context) (firebaseauthenticationinterface.FirebaseAuthentication, error) {
	opt := option.WithCredentialsFile(configs.GetEnvVariable(configs.FirebaseServiceAccountPath))
	config := &firebase.Config{ProjectID: configs.GetEnvVariable(configs.FirebaseConfigProjectID)}
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

func (fbAuth *firebaseAuthentication) VerifyAndDecodeToken(authToken string) (*auth.Token, error) {
	token, err := (*fbAuth).Client.VerifyIDToken(*fbAuth.Context, authToken)
	if err != nil {
		return nil, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.DecodingAuthTokenFailed,
			"/authentication/verifyAndDecodeToken",
			err,
		)
	}
	return token, nil
}

func (fbAuth *firebaseAuthentication) GetAuthUserDataByEmail(email string) (*auth.UserRecord, error) {
	user, err := (*fbAuth).Client.GetUserByEmail(*fbAuth.Context, email)
	if err != nil {
		return nil, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.GetAuthDataFailed,
			"/authentication/getAuthUserDataByEmail",
			err,
		)
	}
	return user, nil
}

func (fbAuth *firebaseAuthentication) GetAuthUserDataById(uid string) (*auth.UserRecord, error) {
	user, err := (*fbAuth).Client.GetUser(*fbAuth.Context, uid)
	if err != nil {
		return nil, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.GetAuthDataFailed,
			"/authentication/getAuthUserDataById",
			err,
		)
	}
	return user, nil
}

func (fbAuth *firebaseAuthentication) SetRoleInAuthUserData(user *auth.UserRecord, accountRole string, dbId string) (bool, error) {
	claims := map[string]interface{}{"type": accountRole, "_id": dbId}
	if err := (*fbAuth).Client.SetCustomUserClaims(*fbAuth.Context, (*user).UID, claims); err != nil {
		return false, horeekaaexception.NewExceptionObject(
			horeekaaexceptionenums.SetAuthDataFailed,
			"/authentication/setRoleInAuthUserData",
			err,
		)
	}

	return true, nil
}

func (fbAuth *firebaseAuthentication) UpdateAuthUserData(user *firebaseauthenticationmodel.UpdateAuthUserData) (*auth.UserRecord, error) {
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
	return updatedUser, nil
}

func (fbAuth *firebaseAuthentication) GenerateEmailVerificationLink(email string) (string, error) {
	link, err := (*fbAuth).Client.EmailVerificationLinkWithSettings(
		*fbAuth.Context,
		email,
		GetFirebaseActionCodeSettings()["data"].(*auth.ActionCodeSettings),
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
		GetFirebaseActionCodeSettings()["data"].(*auth.ActionCodeSettings),
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
