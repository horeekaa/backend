package addressregiongroupdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateAddressRegionGroupTransactionComponent interface {
	PreTransaction(
		createAddressRegionGroupInput *model.InternalCreateAddressRegionGroup,
	) (*model.InternalCreateAddressRegionGroup, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createAddressRegionGroupInput *model.InternalCreateAddressRegionGroup,
	) (*model.AddressRegionGroup, error)

	GenerateNewObjectID() primitive.ObjectID
	GetCurrentObjectID() primitive.ObjectID
}

type CreateAddressRegionGroupRepository interface {
	RunTransaction(
		createAddressRegionGroupInput *model.InternalCreateAddressRegionGroup,
	) (*model.AddressRegionGroup, error)
}
