package addressregiongrouppresentationusecaseinterfaces

import "github.com/horeekaa/backend/model"

type GetAddressRegionGroupUsecase interface {
	Execute(input *model.AddressRegionGroupFilterFields) (*model.AddressRegionGroup, error)
}
