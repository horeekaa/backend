package memberaccessrefpresentationusecaseinterfaces

import (
	memberaccessrefpresentationusecasetypes "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type UpdateMemberAccessRefUsecase interface {
	Execute(input memberaccessrefpresentationusecasetypes.UpdateMemberAccessRefUsecaseInput) (*model.MemberAccessRef, error)
}
