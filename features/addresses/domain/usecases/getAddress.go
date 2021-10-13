package addresspresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	addressdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addresses/domain/repositories"
	addresspresentationusecaseinterfaces "github.com/horeekaa/backend/features/addresses/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getAddressUsecase struct {
	getAddressRepository addressdomainrepositoryinterfaces.GetAddressRepository
}

func NewGetAddressUsecase(
	getAddressRepository addressdomainrepositoryinterfaces.GetAddressRepository,
) (addresspresentationusecaseinterfaces.GetAddressUsecase, error) {
	return &getAddressUsecase{
		getAddressRepository,
	}, nil
}

func (getAddrUsecase *getAddressUsecase) validation(
	input *model.AddressFilterFields,
) (*model.AddressFilterFields, error) {
	return input, nil
}

func (getAddrUsecase *getAddressUsecase) Execute(
	filterFields *model.AddressFilterFields,
) (*model.Address, error) {
	validatedFilterFields, err := getAddrUsecase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	address, err := getAddrUsecase.getAddressRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getAddress",
			err,
		)
	}
	return address, nil
}
