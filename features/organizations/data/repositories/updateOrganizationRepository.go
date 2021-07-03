package organizationdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	organizationdomainrepositorytypes "github.com/horeekaa/backend/features/organizations/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type updateOrganizationRepository struct {
	updateOrganizationTransactionComponent organizationdomainrepositoryinterfaces.UpdateOrganizationTransactionComponent
	mongoDBTransaction                     mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewUpdateOrganizationRepository(
	updateOrganizationRepositoryTransactionComponent organizationdomainrepositoryinterfaces.UpdateOrganizationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (organizationdomainrepositoryinterfaces.UpdateOrganizationRepository, error) {
	updateOrganizationRepo := &updateOrganizationRepository{
		updateOrganizationRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		updateOrganizationRepo,
		"UpdateOrganizationRepository",
	)

	return updateOrganizationRepo, nil
}

func (updateOrgRepo *updateOrganizationRepository) SetValidation(
	usecaseComponent organizationdomainrepositoryinterfaces.UpdateOrganizationUsecaseComponent,
) (bool, error) {
	updateOrgRepo.updateOrganizationTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (updateOrgRepo *updateOrganizationRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateOrgRepo.updateOrganizationTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateOrganization),
	)
}

func (updateOrgRepo *updateOrganizationRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return updateOrgRepo.updateOrganizationTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalUpdateOrganization),
	)
}

func (updateOrgRepo *updateOrganizationRepository) RunTransaction(
	input *model.InternalUpdateOrganization,
) (*organizationdomainrepositorytypes.UpdateOrganizationOutput, error) {
	output, err := updateOrgRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*organizationdomainrepositorytypes.UpdateOrganizationOutput), err
}
