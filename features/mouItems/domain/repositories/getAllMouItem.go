package mouitemdomainrepositoryinterfaces

import (
	mouitemdomainrepositorytypes "github.com/horeekaa/backend/features/mouItems/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAllMouItemRepository interface {
	Execute(filterFields mouitemdomainrepositorytypes.GetAllMouItemInput) ([]*model.MouItem, error)
}
