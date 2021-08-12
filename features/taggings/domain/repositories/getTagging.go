package taggingdomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetTaggingRepository interface {
	Execute(filterFields *model.TaggingFilterFields) (*model.Tagging, error)
}
