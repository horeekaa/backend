package addressregiongrouppresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
	addressregiongrouppresentationusecaseinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getAddressRegionGroupUsecase struct {
	getAddressRegionGroupRepository addressregiongroupdomainrepositoryinterfaces.GetAddressRegionGroupRepository
	pathIdentity                    string
}

func NewGetAddressRegionGroupUsecase(
	getAddressRegionGroupRepository addressregiongroupdomainrepositoryinterfaces.GetAddressRegionGroupRepository,
) (addressregiongrouppresentationusecaseinterfaces.GetAddressRegionGroupUsecase, error) {
	return &getAddressRegionGroupUsecase{
		getAddressRegionGroupRepository,
		"GetAddressRegionGroupUsecase",
	}, nil
}

func (getAddressRegionGroupUcase *getAddressRegionGroupUsecase) validation(
	input *model.AddressRegionGroupFilterFields,
) (*model.AddressRegionGroupFilterFields, error) {
	return input, nil
}

func (getAddressRegionGroupUcase *getAddressRegionGroupUsecase) Execute(
	filterFields *model.AddressRegionGroupFilterFields,
) (*model.AddressRegionGroup, error) {
	validatedFilterFields, err := getAddressRegionGroupUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	addressRegionGroup, err := getAddressRegionGroupUcase.getAddressRegionGroupRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getAddressRegionGroupUcase.pathIdentity,
			err,
		)
	}
	return addressRegionGroup, nil
}
