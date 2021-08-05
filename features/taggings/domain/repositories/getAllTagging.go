package taggingdomainrepositoryinterfaces

import (
	taggingdomainrepositorytypes "github.com/horeekaa/backend/features/taggings/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type GetAllTaggingRepository interface {
	Execute(filterFields taggingdomainrepositorytypes.GetAllTaggingInput) ([]*model.Tagging, error)
}
