package organizationdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	organizationdomainrepositorytypes "github.com/horeekaa/backend/features/organizations/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type UpdateOrganizationUsecaseComponent interface {
	Validation(
		updateOrganizationInput *model.InternalUpdateOrganization,
	) (*model.InternalUpdateOrganization, error)
}

type UpdateOrganizationTransactionComponent interface {
	SetValidation(usecaseComponent UpdateOrganizationUsecaseComponent) (bool, error)

	PreTransaction(
		updateOrganizationInput *model.InternalUpdateOrganization,
	) (*model.InternalUpdateOrganization, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateOrganizationInput *model.InternalUpdateOrganization,
	) (*organizationdomainrepositorytypes.UpdateOrganizationOutput, error)
}

type UpdateOrganizationRepository interface {
	SetValidation(usecaseComponent UpdateOrganizationUsecaseComponent) (bool, error)
	RunTransaction(
		updateOrganizationInput *model.InternalUpdateOrganization,
	) (*organizationdomainrepositorytypes.UpdateOrganizationOutput, error)
}
