package organizationdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateOrganizationUsecaseComponent interface {
	Validation(
		createOrganizationInput *model.InternalCreateOrganization,
	) (*model.InternalCreateOrganization, error)
}

type CreateOrganizationTransactionComponent interface {
	SetValidation(usecaseComponent CreateOrganizationUsecaseComponent) (bool, error)

	PreTransaction(
		createOrganizationInput *model.InternalCreateOrganization,
	) (*model.InternalCreateOrganization, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		createOrganizationInput *model.InternalCreateOrganization,
	) (*model.Organization, error)

	GenerateNewObjectID() primitive.ObjectID
	GetCurrentObjectID() primitive.ObjectID
}

type CreateOrganizationRepository interface {
	SetValidation(usecaseComponent CreateOrganizationUsecaseComponent) (bool, error)
	RunTransaction(
		createOrganizationInput *model.InternalCreateOrganization,
	) (*model.Organization, error)
}
