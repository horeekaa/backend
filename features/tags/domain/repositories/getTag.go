package tagdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetTagRepository interface {
	Execute(filterFields *model.TagFilterFields) (*model.Tag, error)
}
