package descriptivephotodomainrepositoryinterfaces

import "github.com/horeekaa/backend/model"

type GetDescriptivePhotoRepository interface {
	Execute(filterFields *model.DescriptivePhotoFilterFields) (*model.DescriptivePhoto, error)
}
