package accountdomainrepositorytypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

const (
	ManagePersonDeviceTokenActionInsert string = "MANAGE_PERSON_DEVICE_TOKEN_INSERT"
	ManagePersonDeviceTokenActionRemove string = "MANAGE_PERSON_DEVICE_TOKEN_REMOVE"
)

type ManageAccountAuthenticationInput struct {
	AuthHeader string
	Context    context.Context
}

type GetPersonDataByAccountOutput struct {
	Person  *model.Person
	Account *model.Account
}

type ManagePersonDeviceTokenInput struct {
	Person                        *model.Person
	DeviceToken                   string
	ManagePersonDeviceTokenAction string
}
