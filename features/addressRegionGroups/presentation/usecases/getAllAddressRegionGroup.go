package addressregiongrouppresentationusecaseinterfaces

import (
	addressregiongrouppresentationusecasetypes "github.com/horeekaa/backend/features/addressRegionGroups/presentation/usecases/types"
	"github.com/horeekaa/backend/model"
)

type GetAllAddressRegionGroupUsecase interface {
	Execute(input addressregiongrouppresentationusecasetypes.GetAllAddressRegionGroupUsecaseInput) ([]*model.AddressRegionGroup, error)
}
