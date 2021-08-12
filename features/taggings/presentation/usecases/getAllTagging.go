package taggingpresentationusecaseinterfaces

import (
	taggingpresentationusecasetypes "github.com/horeekaa/backend/features/taggings/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAllTaggingUsecase interface {
	Execute(input taggingpresentationusecasetypes.GetAllTaggingUsecaseInput) ([]*model.Tagging, error)
}
