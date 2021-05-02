package accountpresentationusecaseinterfaces

import (
	"github.com/horeekaa/backend/model"
)

type GetAccountUsecase interface {
	Execute(input *model.AccountFilterFields) (*model.Account, error)
}
