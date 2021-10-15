package addressdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type CreateAddressTransactionComponent interface {
	PreTransaction(
		createAddressInput *model.InternalCreateAddress,
	) (*model.InternalCreateAddress, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createAddressInput *model.InternalCreateAddress,
	) (*model.Address, error)
}
