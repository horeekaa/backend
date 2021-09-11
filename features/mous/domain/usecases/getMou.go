package moupresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	moudomainrepositoryinterfaces "github.com/horeekaa/backend/features/mous/domain/repositories"
	moupresentationusecaseinterfaces "github.com/horeekaa/backend/features/mous/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getMouUsecase struct {
	getMouRepository moudomainrepositoryinterfaces.GetMouRepository
}

func NewGetMouUsecase(
	getMouRepository moudomainrepositoryinterfaces.GetMouRepository,
) (moupresentationusecaseinterfaces.GetMouUsecase, error) {
	return &getMouUsecase{
		getMouRepository,
	}, nil
}

func (getOrgUcase *getMouUsecase) validation(
	input *model.MouFilterFields,
) (*model.MouFilterFields, error) {
	return input, nil
}

func (getOrgUcase *getMouUsecase) Execute(
	filterFields *model.MouFilterFields,
) (*model.Mou, error) {
	validatedFilterFields, err := getOrgUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	mou, err := getOrgUcase.getMouRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getMou",
			err,
		)
	}
	return mou, nil
}
