package addressregiongrouppresentationusecaseinterfaces

import (
	addressregiongrouppresentationusecasetypes "github.com/horeekaa/backend/features/addressRegionGroups/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type CreateAddressRegionGroupUsecase interface {
	Execute(input addressregiongrouppresentationusecasetypes.CreateAddressRegionGroupUsecaseInput) (*model.AddressRegionGroup, error)
}
