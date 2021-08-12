package taggingpresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetTaggingUsecase interface {
	Execute(input *model.TaggingFilterFields) (*model.Tagging, error)
}
