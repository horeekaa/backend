package descriptivephotopresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	descriptivephotopresentationusecaseinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getDescriptivePhotoUsecase struct {
	getDescriptivePhotoRepository descriptivephotodomainrepositoryinterfaces.GetDescriptivePhotoRepository
	pathIdentity                  string
}

func NewGetDescriptivePhotoUsecase(
	getDescriptivePhotoRepository descriptivephotodomainrepositoryinterfaces.GetDescriptivePhotoRepository,
) (descriptivephotopresentationusecaseinterfaces.GetDescriptivePhotoUsecase, error) {
	return &getDescriptivePhotoUsecase{
		getDescriptivePhotoRepository,
		"GetDescriptivePhotoUsecase",
	}, nil
}

func (getDescriptivePhotoUcase *getDescriptivePhotoUsecase) validation(
	input *model.DescriptivePhotoFilterFields,
) (*model.DescriptivePhotoFilterFields, error) {
	return input, nil
}

func (getDescriptivePhotoUcase *getDescriptivePhotoUsecase) Execute(
	filterFields *model.DescriptivePhotoFilterFields,
) (*model.DescriptivePhoto, error) {
	validatedFilterFields, err := getDescriptivePhotoUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	descriptivePhoto, err := getDescriptivePhotoUcase.getDescriptivePhotoRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getDescriptivePhotoUcase.pathIdentity,
			err,
		)
	}
	return descriptivePhoto, nil
}
