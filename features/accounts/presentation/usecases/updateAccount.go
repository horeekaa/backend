package accountpresentationusecaseinterfaces

import (
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type UpdateAccountUsecase interface {
	Execute(input accountpresentationusecasetypes.UpdateAccountUsecaseInput) (*model.Account, error)
}
