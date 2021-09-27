package purchaseorderdomainrepositoryutilityinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type PurchaseOrderLoader interface {
	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		mou *model.MouForPurchaseOrderInput,
	) (bool, error)
}
