package authserviceclients

import (
	horeekaaexceptiontofailure "github.com/horeekaa/backend/_errors/serviceFailures/_exceptionToFailure"
	firebaseauthinterfaces "github.com/horeekaa/backend/repositories/authentication/firebase/interfaces"
	authserviceclientinterfaces "github.com/horeekaa/backend/services/authentication/interfaces"
	authservicemodels "github.com/horeekaa/backend/services/authentication/models"
	authserviceoperations "github.com/horeekaa/backend/services/authentication/operations"
)

type authenticationService struct {
	firebaseAuth *firebaseauthinterfaces.FirebaseAuthentication
}

func NewAuthenticationService(firebaseAuth *firebaseauthinterfaces.FirebaseAuthentication) (authserviceclientinterfaces.AuthenticationService, error) {
	return &authenticationService{
		firebaseAuth: firebaseAuth,
	}, nil
}

func (authService *authenticationService) VerifyTokenAndGetUser(authToken string) (chan *authserviceoperations.AuthenticationServiceUserRecord, chan error) {
	authServiceUserRecordChn := make(chan *authserviceoperations.AuthenticationServiceUserRecord)
	errorChn := make(chan error)

	go func() {
		authServiceToken, err := (*authService.firebaseAuth).VerifyAndDecodeToken(authToken)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/authenticationService/verifyToken",
				&err,
			)
			return
		}

		authUserRecord, err := (*authService.firebaseAuth).GetAuthUserDataById(authServiceToken.Token.UID)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/authenticationService/verifyToken",
				&err,
			)
			return
		}

		authServiceUserRecordChn <- &authserviceoperations.AuthenticationServiceUserRecord{
			ServiceUser: authUserRecord,
		}
	}()

	return authServiceUserRecordChn, errorChn
}

func (authService *authenticationService) GetUserFromEmail(email string) (chan *authserviceoperations.AuthenticationServiceUserRecord, chan error) {
	authServiceUserRecordChn := make(chan *authserviceoperations.AuthenticationServiceUserRecord)
	errorChn := make(chan error)

	go func() {
		authUserRecord, err := (*authService.firebaseAuth).GetAuthUserDataByEmail(email)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/authenticationService/getUserFromEmail",
				&err,
			)
			return
		}

		authServiceUserRecordChn <- &authserviceoperations.AuthenticationServiceUserRecord{
			ServiceUser: authUserRecord,
		}
	}()

	return authServiceUserRecordChn, errorChn
}

func (authService *authenticationService) SetRoleInAuthUserData(uid string, accountRole string, dbID string) (chan *authserviceoperations.AuthenticationServiceUserRecord, chan error) {
	authServiceUserRecordChn := make(chan *authserviceoperations.AuthenticationServiceUserRecord)
	errorChn := make(chan error)

	go func() {
		_, err := (*authService.firebaseAuth).SetRoleInAuthUserData(uid, accountRole, dbID)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/authenticationService/setRoleInAuthUserData",
				&err,
			)
			return
		}

		authUserRecord, err := (*authService.firebaseAuth).GetAuthUserDataById(uid)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/authenticationService/setRoleInAuthUserData",
				&err,
			)
			return
		}

		authServiceUserRecordChn <- &authserviceoperations.AuthenticationServiceUserRecord{
			ServiceUser: authUserRecord,
		}
	}()

	return authServiceUserRecordChn, errorChn
}

func (authService *authenticationService) UpdateAuthUserData(updatedUser *authservicemodels.UpdateAuthUserData) (chan *authserviceoperations.AuthenticationServiceUserRecord, chan error) {
	authServiceUserRecordChn := make(chan *authserviceoperations.AuthenticationServiceUserRecord)
	errorChn := make(chan error)

	go func() {
		authUserRecord, err := (*authService.firebaseAuth).UpdateAuthUserData((*updatedUser).ServiceUpdateUser)
		if err != nil {
			errorChn <- horeekaaexceptiontofailure.ConvertException(
				"/authenticationService/updateAuthUserData",
				&err,
			)
			return
		}

		authServiceUserRecordChn <- &authserviceoperations.AuthenticationServiceUserRecord{
			ServiceUser: authUserRecord,
		}
	}()

	return authServiceUserRecordChn, errorChn
}

func (authService *authenticationService) GenerateAuthenticationRelatedLink(email string, action string) (chan string, chan error) {
	urlLinkChn := make(chan string)
	errorChn := make(chan error)

	go func() {
		switch action {
		case authserviceoperations.VerifyEmailAuthURLAction:
			urlLink, err := (*authService.firebaseAuth).GenerateEmailVerificationLink(email)
			if err != nil {
				errorChn <- horeekaaexceptiontofailure.ConvertException(
					"/authenticationService/generateAuthenticationRelatedLink",
					&err,
				)
				return
			}
			urlLinkChn <- urlLink
			break

		case authserviceoperations.ResetPasswordAuthURLAction:
			urlLink, err := (*authService.firebaseAuth).GeneratePasswordResetLink(email)
			if err != nil {
				errorChn <- horeekaaexceptiontofailure.ConvertException(
					"/authenticationService/generateAuthenticationRelatedLink",
					&err,
				)
				return
			}
			urlLinkChn <- urlLink
			break
		}
	}()

	return urlLinkChn, errorChn
}
