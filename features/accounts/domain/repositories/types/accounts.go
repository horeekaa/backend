package accountdomainrepositorytypes

import (
	"context"

	"firebase.google.com/go/v4/auth"
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

type ManageAccountAuthenticationInput struct {
	AuthHeader string
	Context    context.Context
}

type GetUserFromAuthHeaderInput struct {
	AuthHeader string
	Context    context.Context
}

type GetAccountFromAuthDataInput struct {
	Context context.Context
	User    *auth.UserRecord
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
