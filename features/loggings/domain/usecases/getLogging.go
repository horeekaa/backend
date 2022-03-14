package loggingpresentationusecases

import (
	horeekaacorefailuretoerror "github.com/horeekaa/backend/core/errors/errors/failureToError"
	loggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/loggings/domain/repositories"
	loggingpresentationusecaseinterfaces "github.com/horeekaa/backend/features/loggings/presentation/usecases"
	"github.com/horeekaa/backend/model"
)

type getLoggingUsecase struct {
	getLoggingRepository loggingdomainrepositoryinterfaces.GetLoggingRepository
	pathIdentity         string
}

func NewGetLoggingUsecase(
	getLoggingRepository loggingdomainrepositoryinterfaces.GetLoggingRepository,
) (loggingpresentationusecaseinterfaces.GetLoggingUsecase, error) {
	return &getLoggingUsecase{
		getLoggingRepository,
		"GetLoggingUsecase",
	}, nil
}

func (getLogUcase *getLoggingUsecase) validation(
	input *model.LoggingFilterFields,
) (*model.LoggingFilterFields, error) {
	return input, nil
}

func (getLogUcase *getLoggingUsecase) Execute(
	filterFields *model.LoggingFilterFields,
) (*model.Logging, error) {
	validatedFilterFields, err := getLogUcase.validation(filterFields)
	if err != nil {
		return nil, err
	}

	logging, err := getLogUcase.getLoggingRepository.Execute(
		validatedFilterFields,
	)
	if err != nil {
		return nil, horeekaacorefailuretoerror.ConvertFailure(
			getLogUcase.pathIdentity,
			err,
		)
	}
	return logging, nil
}
