package memberaccesspresentationusecaseinterfaces

import (
	memberaccesspresentationusecasetypes "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type CreateMemberAccessUsecase interface {
	Execute(input memberaccesspresentationusecasetypes.CreateMemberAccessUsecaseInput) (*model.MemberAccess, error)
}
