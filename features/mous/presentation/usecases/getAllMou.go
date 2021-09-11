package moupresentationusecaseinterfaces

import (
	moupresentationusecasetypes "github.com/horeekaa/backend/features/mous/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAllMouUsecase interface {
	Execute(input moupresentationusecasetypes.GetAllMouUsecaseInput) ([]*model.Mou, error)
}
