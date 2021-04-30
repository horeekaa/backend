package memberaccessrefpresentationusecaseinterfaces

import (
	memberaccessrefpresentationusecasetypes "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAllMemberAccessRefUsecase interface {
	Execute(input memberaccessrefpresentationusecasetypes.GetAllMemberAccessRefUsecaseInput) ([]*model.MemberAccessRef, error)
}
