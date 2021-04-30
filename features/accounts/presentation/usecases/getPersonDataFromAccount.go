package accountpresentationusecaseinterfaces

import (
	accountpresentationusecasetypes "github.com/horeekaa/backend/features/accounts/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetPersonDataFromAccountUsecase interface {
	Execute(input accountpresentationusecasetypes.GetPersonDataFromAccountInput) (*model.Person, error)
}
