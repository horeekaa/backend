package addressregiongrouppresentationusecasetypes

import (
	"context"

	"github.com/horeekaa/backend/model"
)

type CreateAddressRegionGroupUsecaseInput struct {
	Context   context.Context
	CreateAddressRegionGroup *model.CreateAddressRegionGroup
}

type UpdateAddressRegionGroupUsecaseInput struct {
	Context   context.Context
	UpdateAddressRegionGroup *model.UpdateAddressRegionGroup
}

type GetAllAddressRegionGroupUsecaseInput struct {
	Context       context.Context
	FilterFields  *model.AddressRegionGroupFilterFields
	PaginationOps *model.PaginationOptionInput
}
