package addressregiongrouppresentationusecaseinterfaces

import (
	addressregiongrouppresentationusecasetypes "github.com/horeekaa/backend/features/addressRegionGroups/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type UpdateAddressRegionGroupUsecase interface {
	Execute(input addressregiongrouppresentationusecasetypes.UpdateAddressRegionGroupUsecaseInput) (*model.AddressRegionGroup, error)
}
