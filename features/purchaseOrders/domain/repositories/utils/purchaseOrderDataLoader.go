package purchaseorderdomainrepositoryutilityinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type PurchaseOrderLoader interface {
	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		mou *model.MouForPurchaseOrderInput,
		organization *model.OrganizationForPurchaseOrderInput,
	) (bool, error)
}
