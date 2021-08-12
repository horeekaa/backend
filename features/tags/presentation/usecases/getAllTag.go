package tagpresentationusecaseinterfaces

import (
	tagpresentationusecasetypes "github.com/horeekaa/backend/features/tags/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAllTagUsecase interface {
	Execute(input tagpresentationusecasetypes.GetAllTagUsecaseInput) ([]*model.Tag, error)
}
