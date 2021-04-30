package organizationdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
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
	mongoDBTransaction.SetTransaction(
		updateOrganizationRepositoryTransactionComponent,
		"UpdateOrganizationRepository",
	)

	return &updateOrganizationRepository{
		updateOrganizationTransactionComponent: updateOrganizationRepositoryTransactionComponent,
		mongoDBTransaction:                     mongoDBTransaction,
	}, nil
}

func (updateMmbAccRefRepo *updateOrganizationRepository) SetValidation(
	usecaseComponent organizationdomainrepositoryinterfaces.UpdateOrganizationUsecaseComponent,
) (bool, error) {
	updateMmbAccRefRepo.updateOrganizationTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (updateMmbAccRefRepo *updateOrganizationRepository) RunTransaction(
	input *model.UpdateOrganization,
) (*organizationdomainrepositorytypes.UpdateOrganizationOutput, error) {
	output, err := updateMmbAccRefRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*organizationdomainrepositorytypes.UpdateOrganizationOutput), err
}
