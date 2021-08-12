package taggingpresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	taggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/taggings/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getTaggingUsecase struct {
	getTaggingRepository taggingdomainrepositoryinterfaces.GetTaggingRepository
}

func NewGetTaggingUsecase(
	getTaggingRepository taggingdomainrepositoryinterfaces.GetTaggingRepository,
) (taggingpresentationusecaseinterfaces.GetTaggingUsecase, error) {
	return &getTaggingUsecase{
		getTaggingRepository,
	}, nil
}

func (getTaggingUcase *getTaggingUsecase) validation(
	input *model.TaggingFilterFields,
) (*model.TaggingFilterFields, error) {
	return input, nil
}

func (getTaggingUcase *getTaggingUsecase) Execute(
	filterFields *model.TaggingFilterFields,
) (*model.Tagging, error) {
	validatedFilterFields, err := getTaggingUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	tagging, err := getTaggingUcase.getTaggingRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getTagging",
			err,
		)
	}
	return tagging, nil
}
