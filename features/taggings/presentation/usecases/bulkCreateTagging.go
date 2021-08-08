package taggingpresentationusecaseinterfaces

import (
	taggingpresentationusecasetypes "github.com/horeekaa/backend/features/taggings/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type BulkCreateTaggingUsecase interface {
	Execute(input taggingpresentationusecasetypes.BulkCreateTaggingUsecaseInput) ([]*model.Tagging, error)
}
