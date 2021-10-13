package addresspresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetAddressUsecase interface {
	Execute(input *model.AddressFilterFields) (*model.Address, error)
}
