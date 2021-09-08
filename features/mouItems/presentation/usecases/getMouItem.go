package mouitempresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetMouItemUsecase interface {
	Execute(input *model.MouItemFilterFields) (*model.MouItem, error)
}
