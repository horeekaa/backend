package descriptivephotopresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	descriptivephotopresentationusecaseinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getDescriptivePhotoUsecase struct {
	getDescriptivePhotoRepository descriptivephotodomainrepositoryinterfaces.GetDescriptivePhotoRepository
}

func NewGetDescriptivePhotoUsecase(
	getDescriptivePhotoRepository descriptivephotodomainrepositoryinterfaces.GetDescriptivePhotoRepository,
) (descriptivephotopresentationusecaseinterfaces.GetDescriptivePhotoUsecase, error) {
	return &getDescriptivePhotoUsecase{
		getDescriptivePhotoRepository,
	}, nil
}

func (getLogUcase *getDescriptivePhotoUsecase) validation(
	input *model.DescriptivePhotoFilterFields,
) (*model.DescriptivePhotoFilterFields, error) {
	return input, nil
}

func (getLogUcase *getDescriptivePhotoUsecase) Execute(
	filterFields *model.DescriptivePhotoFilterFields,
) (*model.DescriptivePhoto, error) {
	validatedFilterFields, err := getLogUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	descriptivePhoto, err := getLogUcase.getDescriptivePhotoRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getDescriptivePhoto",
			err,
		)
	}
	return descriptivePhoto, nil
}
