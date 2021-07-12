package organizationdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ProposeUpdateOrganizationUsecaseComponent interface {
	Validation(
		updateOrganizationInput *model.InternalUpdateOrganization,
	) (*model.InternalUpdateOrganization, error)
}

type ProposeUpdateOrganizationTransactionComponent interface {
	SetValidation(usecaseComponent ProposeUpdateOrganizationUsecaseComponent) (bool, error)

	PreTransaction(
		updateOrganizationInput *model.InternalUpdateOrganization,
	) (*model.InternalUpdateOrganization, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateOrganizationInput *model.InternalUpdateOrganization,
	) (*model.Organization, error)
}

type ProposeUpdateOrganizationRepository interface {
	SetValidation(usecaseComponent ProposeUpdateOrganizationUsecaseComponent) (bool, error)
	RunTransaction(
		updateOrganizationInput *model.InternalUpdateOrganization,
	) (*model.Organization, error)
}
