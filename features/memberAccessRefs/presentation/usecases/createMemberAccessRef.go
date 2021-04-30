package memberaccessrefpresentationusecaseinterfaces

import (
	memberaccessrefpresentationusecasetypes "github.com/horeekaa/backend/features/memberAccessRefs/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type CreateMemberAccessRefUsecase interface {
	Execute(input memberaccessrefpresentationusecasetypes.CreateMemberAccessRefUsecaseInput) (*model.MemberAccessRef, error)
}
