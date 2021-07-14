package organizationdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateOrganizationRepository struct {
	approveUpdateOrganizationTransactionComponent organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationTransactionComponent
	mongoDBTransaction                            mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewApproveUpdateOrganizationRepository(
	approveUpdateOrganizationRepositoryTransactionComponent organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationRepository, error) {
	approveUpdateOrganizationRepo := &approveUpdateOrganizationRepository{
		approveUpdateOrganizationRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		approveUpdateOrganizationRepo,
		"ApproveUpdateOrganizationRepository",
	)

	return approveUpdateOrganizationRepo, nil
}

func (approveUpdateOrgRepo *approveUpdateOrganizationRepository) SetValidation(
	usecaseComponent organizationdomainrepositoryinterfaces.ApproveUpdateOrganizationUsecaseComponent,
) (bool, error) {
	approveUpdateOrgRepo.approveUpdateOrganizationTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (approveUpdateOrgRepo *approveUpdateOrganizationRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return approveUpdateOrgRepo.approveUpdateOrganizationTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateOrganization),
	)
}

func (approveUpdateOrgRepo *approveUpdateOrganizationRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return approveUpdateOrgRepo.approveUpdateOrganizationTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalUpdateOrganization),
	)
}

func (approveUpdateOrgRepo *approveUpdateOrganizationRepository) RunTransaction(
	input *model.InternalUpdateOrganization,
) (*model.Organization, error) {
	output, err := approveUpdateOrgRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.Organization), err
}
