package moupresentationusecaseinterfaces

import (
	moupresentationusecasetypes "github.com/horeekaa/backend/features/mous/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type UpdateMouUsecase interface {
	Execute(input moupresentationusecasetypes.UpdateMouUsecaseInput) (*model.Mou, error)
}
