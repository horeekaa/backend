package servicecoordinatormodels

import "github.com/horeekaa/backend/model"

const (
	ManagePersonDeviceTokenActionInsert string = "MANAGE_PERSON_DEVICE_TOKEN_INSERT"
	ManagePersonDeviceTokenActionRemove string = "MANAGE_PERSON_DEVICE_TOKEN_REMOVE"
)

type GetPersonDataByAccountOutput struct {
	Person  *model.Person
	Account *model.Account
}

type ManagePersonDeviceTokenInput struct {
	Person                        *model.Person
	DeviceToken                   string
	ManagePersonDeviceTokenAction string
}
