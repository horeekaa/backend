package mouitemdomainrepositoryutilityinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type AgreedProductLoader interface {
	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		product *model.ObjectIDOnly,
		agreedProduct *model.InternalAgreedProductInput,
		organization *model.OrganizationForMouItemInput,
	) (bool, error)
}
