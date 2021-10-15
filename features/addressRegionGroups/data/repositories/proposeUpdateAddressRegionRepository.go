package addressregiongroupdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateAddressRegionGroupRepository struct {
	proposeUpdateAddressRegionGroupTransactionComponent addressregiongroupdomainrepositoryinterfaces.ProposeUpdateAddressRegionGroupTransactionComponent
	mongoDBTransaction                                  mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewProposeUpdateAddressRegionGroupRepository(
	proposeUpdateAddressRegionGroupRepositoryTransactionComponent addressregiongroupdomainrepositoryinterfaces.ProposeUpdateAddressRegionGroupTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (addressregiongroupdomainrepositoryinterfaces.ProposeUpdateAddressRegionGroupRepository, error) {
	proposeUpdateAddressRegionGroupRepo := &proposeUpdateAddressRegionGroupRepository{
		proposeUpdateAddressRegionGroupRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		proposeUpdateAddressRegionGroupRepo,
		"ProposeUpdateAddressRegionGroupRepository",
	)

	return proposeUpdateAddressRegionGroupRepo, nil
}

func (updateAddressRegionGroupRepo *proposeUpdateAddressRegionGroupRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateAddressRegionGroupRepo.proposeUpdateAddressRegionGroupTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateAddressRegionGroup),
	)
}

func (updateAddressRegionGroupRepo *proposeUpdateAddressRegionGroupRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	addressRegionGroupToUpdate := input.(*model.InternalUpdateAddressRegionGroup)

	return updateAddressRegionGroupRepo.proposeUpdateAddressRegionGroupTransactionComponent.TransactionBody(
		operationOption,
		addressRegionGroupToUpdate,
	)
}

func (updateAddressRegionGroupRepo *proposeUpdateAddressRegionGroupRepository) RunTransaction(
	input *model.InternalUpdateAddressRegionGroup,
) (*model.AddressRegionGroup, error) {
	output, err := updateAddressRegionGroupRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.AddressRegionGroup), err
}
