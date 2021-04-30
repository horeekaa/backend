package organizationdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	organizationdomainrepositorytypes "github.com/horeekaa/backend/features/organizations/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type UpdateOrganizationUsecaseComponent interface {
	Validation(
		updateOrganizationInput *model.UpdateOrganization,
	) (*model.UpdateOrganization, error)
}

type UpdateOrganizationTransactionComponent interface {
	SetValidation(usecaseComponent UpdateOrganizationUsecaseComponent) (bool, error)

	PreTransaction(
		updateOrganizationInput *model.UpdateOrganization,
	) (*model.UpdateOrganization, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateOrganizationInput *model.UpdateOrganization,
	) (*organizationdomainrepositorytypes.UpdateOrganizationOutput, error)
}

type UpdateOrganizationRepository interface {
	SetValidation(usecaseComponent UpdateOrganizationUsecaseComponent) (bool, error)
	RunTransaction(
		updateOrganizationInput *model.UpdateOrganization,
	) (*organizationdomainrepositorytypes.UpdateOrganizationOutput, error)
}
