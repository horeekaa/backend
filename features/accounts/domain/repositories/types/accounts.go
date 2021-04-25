package accountdomainrepositorytypes

import (
	"context"

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
