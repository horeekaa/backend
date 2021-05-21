package accountdomainrepositorytypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

const (
	ManageAccountDeviceTokenActionInsert string = "MANAGE_ACCOUNT_DEVICE_TOKEN_INSERT"
	ManageAccountDeviceTokenActionRemove string = "MANAGE_ACCOUNT_DEVICE_TOKEN_REMOVE"
)

type CreateAccountFromAuthDataInput struct {
	Context context.Context
}

type GetUserFromAuthHeaderInput struct {
	AuthHeader string
	Context    context.Context
}

type GetAccountFromAuthDataInput struct {
	Context context.Context
}

type ManageAccountDeviceTokenInput struct {
	Account                        *model.Account
	DeviceToken                    string
	ManageAccountDeviceTokenAction string
}
