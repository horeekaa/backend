package authserviceclientinterfaces

import (
	authservicemodels "github.com/horeekaa/backend/services/authentication/models"
	authserviceoperations "github.com/horeekaa/backend/services/authentication/operations"
)

type AuthenticationService interface {
	VerifyTokenAndGetUser(authToken string) (chan *authserviceoperations.AuthenticationServiceUserRecord, chan error)
	GetUserFromEmail(email string) (chan *authserviceoperations.AuthenticationServiceUserRecord, chan error)
	SetRoleInAuthUserData(uid string, accountType string, dbID string) (chan *authserviceoperations.AuthenticationServiceUserRecord, chan error)
	UpdateAuthUserData(updatedUser *authservicemodels.UpdateAuthUserData) (chan *authserviceoperations.AuthenticationServiceUserRecord, chan error)
	GenerateAuthenticationRelatedLink(email string, action string) (chan string, chan error)
}
