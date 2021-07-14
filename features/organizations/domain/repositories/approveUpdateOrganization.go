package organizationdomainrepositoryinterfaces

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	"github.com/horeekaa/backend/model"
)

type ApproveUpdateOrganizationUsecaseComponent interface {
	Validation(
		updateOrganizationInput *model.InternalUpdateOrganization,
	) (*model.InternalUpdateOrganization, error)
}

type ApproveUpdateOrganizationTransactionComponent interface {
	SetValidation(usecaseComponent ApproveUpdateOrganizationUsecaseComponent) (bool, error)

	PreTransaction(
		updateOrganizationInput *model.InternalUpdateOrganization,
	) (*model.InternalUpdateOrganization, error)

	TransactionBody(
		session *mongodbcoretypes.OperationOptions,
		updateOrganizationInput *model.InternalUpdateOrganization,
	) (*model.Organization, error)
}

type ApproveUpdateOrganizationRepository interface {
	SetValidation(usecaseComponent ApproveUpdateOrganizationUsecaseComponent) (bool, error)
	RunTransaction(
		updateOrganizationInput *model.InternalUpdateOrganization,
	) (*model.Organization, error)
}
