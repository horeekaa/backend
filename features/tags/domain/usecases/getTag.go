package tagpresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	tagpresentationusecaseinterfaces "github.com/horeekaa/backend/features/tags/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getTagUsecase struct {
	getTagRepository tagdomainrepositoryinterfaces.GetTagRepository
}

func NewGetTagUsecase(
	getTagRepository tagdomainrepositoryinterfaces.GetTagRepository,
) (tagpresentationusecaseinterfaces.GetTagUsecase, error) {
	return &getTagUsecase{
		getTagRepository,
	}, nil
}

func (getTagUcase *getTagUsecase) validation(
	input *model.TagFilterFields,
) (*model.TagFilterFields, error) {
	return input, nil
}

func (getTagUcase *getTagUsecase) Execute(
	filterFields *model.TagFilterFields,
) (*model.Tag, error) {
	validatedFilterFields, err := getTagUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	tag, err := getTagUcase.getTagRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			"/getTag",
			err,
		)
	}
	return tag, nil
}
