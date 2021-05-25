package memberaccesspresentationusecaseinterfaces

import (
	memberaccesspresentationusecasetypes "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type UpdateMemberAccessUsecase interface {
	Execute(input memberaccesspresentationusecasetypes.UpdateMemberAccessUsecaseInput) (*model.MemberAccess, error)
}
