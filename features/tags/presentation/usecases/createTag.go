package tagpresentationusecaseinterfaces

import (
	tagpresentationusecasetypes "github.com/horeekaa/backend/features/tags/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type CreateTagUsecase interface {
	Execute(input tagpresentationusecasetypes.CreateTagUsecaseInput) (*model.Tag, error)
}
