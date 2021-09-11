package moupresentationusecaseinterfaces

import (
	moupresentationusecasetypes "github.com/horeekaa/backend/features/mous/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type CreateMouUsecase interface {
	Execute(input moupresentationusecasetypes.CreateMouUsecaseInput) (*model.Mou, error)
}
