package accountpresentationusecaseinterfaces

import (
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type LogoutUsecase interface {
	Execute(input accountpresentationusecasetypes.LogoutUsecaseInput) (*model.Account, error)
}
