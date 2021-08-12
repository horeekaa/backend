package tagpresentationusecaseinterfaces

import (
	tagpresentationusecasetypes "github.com/horeekaa/backend/features/tags/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type UpdateTagUsecase interface {
	Execute(input tagpresentationusecasetypes.UpdateTagUsecaseInput) (*model.Tag, error)
}
