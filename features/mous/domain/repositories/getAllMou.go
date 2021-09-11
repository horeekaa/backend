package moudomainrepositoryinterfaces

import (
	moudomainrepositorytypes "github.com/horeekaa/backend/features/mous/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAllMouRepository interface {
	Execute(filterFields moudomainrepositorytypes.GetAllMouInput) ([]*model.Mou, error)
}
