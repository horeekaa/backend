package organizationdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	organizationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/organizations/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createOrganizationRepository struct {
	createOrganizationTransactionComponent organizationdomainrepositoryinterfaces.CreateOrganizationTransactionComponent
	mongoDBTransaction                     mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreateOrganizationRepository(
	createOrganizationRepositoryTransactionComponent organizationdomainrepositoryinterfaces.CreateOrganizationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (organizationdomainrepositoryinterfaces.CreateOrganizationRepository, error) {
	createOrganizationRepo := &createOrganizationRepository{
		createOrganizationRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		createOrganizationRepo,
		"CreateOrganizationRepository",
	)

	return createOrganizationRepo, nil
}

func (createOrgRepo *createOrganizationRepository) SetValidation(
	usecaseComponent organizationdomainrepositoryinterfaces.CreateOrganizationUsecaseComponent,
) (bool, error) {
	createOrgRepo.createOrganizationTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (createOrgRepo *createOrganizationRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return createOrgRepo.createOrganizationTransactionComponent.PreTransaction(
		input.(*model.InternalCreateOrganization),
	)
}

func (createOrgRepo *createOrganizationRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return createOrgRepo.createOrganizationTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalCreateOrganization),
	)
}

func (createOrgRepo *createOrganizationRepository) RunTransaction(
	input *model.InternalCreateOrganization,
) (*model.Organization, error) {
	output, err := createOrgRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.Organization), err
}
