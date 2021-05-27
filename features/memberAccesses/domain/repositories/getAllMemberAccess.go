package memberaccessdomainrepositoryinterfaces

import (
	memberaccessdomainrepositorytypes "github.com/horeekaa/backend/features/memberAccesses/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAllMemberAccessRepository interface {
	Execute(filterFields memberaccessdomainrepositorytypes.GetAllMemberAccessInput) ([]*model.MemberAccess, error)
}
