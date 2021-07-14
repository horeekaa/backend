package organizationdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateOrganizationRepository struct {
	proposeUpdateOrganizationTransactionComponent organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationTransactionComponent
	mongoDBTransaction                            mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewProposeUpdateOrganizationRepository(
	proposeUpdateOrganizationRepositoryTransactionComponent organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationRepository, error) {
	proposeUpdateOrganizationRepo := &proposeUpdateOrganizationRepository{
		proposeUpdateOrganizationRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		proposeUpdateOrganizationRepo,
		"ProposeUpdateOrganizationRepository",
	)

	return proposeUpdateOrganizationRepo, nil
}

func (updateOrgRepo *proposeUpdateOrganizationRepository) SetValidation(
	usecaseComponent organizationdomainrepositoryinterfaces.ProposeUpdateOrganizationUsecaseComponent,
) (bool, error) {
	updateOrgRepo.proposeUpdateOrganizationTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (updateOrgRepo *proposeUpdateOrganizationRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateOrgRepo.proposeUpdateOrganizationTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateOrganization),
	)
}

func (updateOrgRepo *proposeUpdateOrganizationRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return updateOrgRepo.proposeUpdateOrganizationTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalUpdateOrganization),
	)
}

func (updateOrgRepo *proposeUpdateOrganizationRepository) RunTransaction(
	input *model.InternalUpdateOrganization,
) (*model.Organization, error) {
	output, err := updateOrgRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.Organization), err
}
