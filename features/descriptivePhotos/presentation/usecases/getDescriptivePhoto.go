package descriptivephotopresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetDescriptivePhotoUsecase interface {
	Execute(input *model.DescriptivePhotoFilterFields) (*model.DescriptivePhoto, error)
}
