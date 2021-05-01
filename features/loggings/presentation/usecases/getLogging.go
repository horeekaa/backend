package loggingpresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetLoggingUsecase interface {
	Execute(input *model.LoggingFilterFields) (*model.Logging, error)
}
