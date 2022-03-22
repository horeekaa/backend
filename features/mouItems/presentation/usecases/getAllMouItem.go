package mouitempresentationusecaseinterfaces

import (
	mouitempresentationusecasetypes "github.com/horeekaa/backend/features/mouItems/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAllMouItemUsecase interface {
	Execute(input mouitempresentationusecasetypes.GetAllMouItemUsecaseInput) ([]*model.MouItem, error)
}
