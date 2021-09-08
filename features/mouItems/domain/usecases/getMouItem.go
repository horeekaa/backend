package mouitempresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	mouitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/mouItems/domain/repositories"
	mouitempresentationusecaseinterfaces "github.com/horeekaa/backend/features/mouItems/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getMouItemUsecase struct {
	getMouItemRepository mouitemdomainrepositoryinterfaces.GetMouItemRepository
}

func NewGetMouItemUsecase(
	getMouItemRepository mouitemdomainrepositoryinterfaces.GetMouItemRepository,
) (mouitempresentationusecaseinterfaces.GetMouItemUsecase, error) {
	return &getMouItemUsecase{
		getMouItemRepository,
	}, nil
}

func (getMouItemUcase *getMouItemUsecase) validation(
	input *model.MouItemFilterFields,
) (*model.MouItemFilterFields, error) {
	return input, nil
}

func (getMouItemUcase *getMouItemUsecase) Execute(
	filterFields *model.MouItemFilterFields,
) (*model.MouItem, error) {
	validatedFilterFields, err := getMouItemUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	mouItem, err := getMouItemUcase.getMouItemRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getMouItem",
			err,
		)
	}
	return mouItem, nil
}
