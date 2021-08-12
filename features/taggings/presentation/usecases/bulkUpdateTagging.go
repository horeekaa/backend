package taggingpresentationusecaseinterfaces

import (
	taggingpresentationusecasetypes "github.com/horeekaa/backend/features/taggings/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type BulkUpdateTaggingUsecase interface {
	Execute(input taggingpresentationusecasetypes.BulkUpdateTaggingUsecaseInput) ([]*model.Tagging, error)
}
