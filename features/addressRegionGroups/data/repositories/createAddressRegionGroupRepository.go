package addressregiongroupdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createAddressRegionGroupRepository struct {
	createAddressRegionGroupTransactionComponent addressregiongroupdomainrepositoryinterfaces.CreateAddressRegionGroupTransactionComponent
	mongoDBTransaction                           mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreateAddressRegionGroupRepository(
	createAddressRegionGroupRepositoryTransactionComponent addressregiongroupdomainrepositoryinterfaces.CreateAddressRegionGroupTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (addressregiongroupdomainrepositoryinterfaces.CreateAddressRegionGroupRepository, error) {
	createAddressRegionGroupRepo := &createAddressRegionGroupRepository{
		createAddressRegionGroupRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		createAddressRegionGroupRepo,
		"CreateAddressRegionGroupRepository",
	)

	return createAddressRegionGroupRepo, nil
}

func (createAddressRegionGroupRepo *createAddressRegionGroupRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return createAddressRegionGroupRepo.createAddressRegionGroupTransactionComponent.PreTransaction(
		input.(*model.InternalCreateAddressRegionGroup),
	)
}

func (createAddressRegionGroupRepo *createAddressRegionGroupRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	addressRegionGroupToCreate := input.(*model.InternalCreateAddressRegionGroup)

	return createAddressRegionGroupRepo.createAddressRegionGroupTransactionComponent.TransactionBody(
		operationOption,
		addressRegionGroupToCreate,
	)
}

func (createAddressRegionGroupRepo *createAddressRegionGroupRepository) RunTransaction(
	input *model.InternalCreateAddressRegionGroup,
) (*model.AddressRegionGroup, error) {
	output, err := createAddressRegionGroupRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.AddressRegionGroup), err
}
