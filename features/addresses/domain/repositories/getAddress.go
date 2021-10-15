package addressdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetAddressRepository interface {
	Execute(filterFields *model.AddressFilterFields) (*model.Address, error)
}
