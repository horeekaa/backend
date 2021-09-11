package mouitemdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetMouItemRepository interface {
	Execute(filterFields *model.MouItemFilterFields) (*model.MouItem, error)
}
