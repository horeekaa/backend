package accountdomainrepositorytypes

import (
	"context"

	authenticationcoremodels "github.com/horeekaa/backend/core/authentication/models"
	"github.com/horeekaa/backend/model"
)

const (
	ManageAccountDeviceTokenActionInsert string = "MANAGE_ACCOUNT_DEVICE_TOKEN_INSERT"
	ManageAccountDeviceTokenActionRemove string = "MANAGE_ACCOUNT_DEVICE_TOKEN_REMOVE"
)

type CreateMemberAccessForAccountInput struct {
	Account                    *model.Account
	MemberAccessRefType        model.MemberAccessRefType
	OrganizationMembershipRole model.OrganizationMembershipRole
	OrganizationType           model.OrganizationType
	Organization               *model.Organization
}

type CreateAccountFromAuthDataInput struct {
	User *authenticationcoremodels.AuthUserWrap
}

type GetUserFromAuthHeaderInput struct {
	AuthHeader string
	Context    context.Context
}

type GetAccountFromAuthDataInput struct {
	Context context.Context
	User    *authenticationcoremodels.AuthUserWrap
}

type ManageAccountDeviceTokenInput struct {
	Account                        *model.Account
	DeviceToken                    string
	ManageAccountDeviceTokenAction string
}

type GetAccountMemberAccessInput struct {
	Account                *model.Account
	MemberAccessRefType    model.MemberAccessRefType
	MemberAccessRefOptions model.MemberAccessRefOptionsInput
}
