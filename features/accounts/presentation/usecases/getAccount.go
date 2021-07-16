package accountpresentationusecaseinterfaces

import (
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAccountUsecase interface {
	Execute(input accountpresentationusecasetypes.GetAccountInput) (*model.Account, error)
}
