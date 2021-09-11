package moupresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetMouUsecase interface {
	Execute(input *model.MouFilterFields) (*model.Mou, error)
}
