package memberaccessrefdomainrepositoryinterfaces

import (
	memberaccessrefdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccessRefs/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAllMemberAccessRefRepository interface {
	Execute(filterFields memberaccessrefdomainrepositorytypes.GetAllMemberAccessRefInput) ([]*model.MemberAccessRef, error)
}
