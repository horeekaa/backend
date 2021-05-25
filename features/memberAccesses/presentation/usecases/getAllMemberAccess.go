package memberaccesspresentationusecaseinterfaces

import (
	memberaccesspresentationusecasetypes "github.com/horeekaa/backend/features/memberAccesses/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAllMemberAccessUsecase interface {
	Execute(input memberaccesspresentationusecasetypes.GetAllMemberAccessUsecaseInput) ([]*model.MemberAccess, error)
}
