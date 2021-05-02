package loggingdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetLoggingRepository interface {
	Execute(filterFields *model.LoggingFilterFields) (*model.Logging, error)
}
