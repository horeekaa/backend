package accountpresentationusecaseinterfaces

import (
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type LoginUsecase interface {
	Execute(input accountpresentationusecasetypes.LoginUsecaseInput) (*model.Account, error)
}
