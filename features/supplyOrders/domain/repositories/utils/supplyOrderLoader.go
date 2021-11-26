package supplyorderdomainrepositoryutilityinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type SupplyOrderLoader interface {
	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		organization *model.OrganizationForSupplyOrderInput,
	) (bool, error)
}
