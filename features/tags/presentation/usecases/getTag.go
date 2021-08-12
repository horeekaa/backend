package tagpresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetTagUsecase interface {
	Execute(input *model.TagFilterFields) (*model.Tag, error)
}
