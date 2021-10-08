package addressdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type UpdateAddressTransactionComponent interface {
	PreTransaction(
		updateAddressInput *model.InternalUpdateAddress,
	) (*model.InternalUpdateAddress, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateAddressInput *model.InternalUpdateAddress,
	) (*model.Address, error)
}
