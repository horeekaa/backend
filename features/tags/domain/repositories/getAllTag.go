package tagdomainrepositoryinterfaces

import (
	tagdomainrepositorytypes "github.com/horeekaa/backend/features/tags/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAllTagRepository interface {
	Execute(filterFields tagdomainrepositorytypes.GetAllTagInput) ([]*model.Tag, error)
}
