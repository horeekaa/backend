package addressregiongroupdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ProposeUpdateAddressRegionGroupTransactionComponent interface {
	PreTransaction(
		updateAddressRegionGroupInput *model.InternalUpdateAddressRegionGroup,
	) (*model.InternalUpdateAddressRegionGroup, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateAddressRegionGroupInput *model.InternalUpdateAddressRegionGroup,
	) (*model.AddressRegionGroup, error)
}

type ProposeUpdateAddressRegionGroupRepository interface {
	RunTransaction(
		updateAddressRegionGroupInput *model.InternalUpdateAddressRegionGroup,
	) (*model.AddressRegionGroup, error)
}
