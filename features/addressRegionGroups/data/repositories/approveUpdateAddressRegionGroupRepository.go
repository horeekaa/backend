package addressregiongroupdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	addressregiongroupdomainrepositoryinterfaces "github.com/horeekaa/backend/features/addressRegionGroups/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateAddressRegionGroupRepository struct {
	approveUpdateAddressRegionGroupTransactionComponent addressregiongroupdomainrepositoryinterfaces.ApproveUpdateAddressRegionGroupTransactionComponent
	mongoDBTransaction                                  mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewApproveUpdateAddressRegionGroupRepository(
	approveUpdateAddressRegionGroupRepositoryTransactionComponent addressregiongroupdomainrepositoryinterfaces.ApproveUpdateAddressRegionGroupTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (addressregiongroupdomainrepositoryinterfaces.ApproveUpdateAddressRegionGroupRepository, error) {
	approveUpdateAddressRegionGroupRepo := &approveUpdateAddressRegionGroupRepository{
		approveUpdateAddressRegionGroupRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		approveUpdateAddressRegionGroupRepo,
		"ApproveUpdateAddressRegionGroupRepository",
	)

	return approveUpdateAddressRegionGroupRepo, nil
}

func (updateAddressRegionGroupRepo *approveUpdateAddressRegionGroupRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateAddressRegionGroupRepo.approveUpdateAddressRegionGroupTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateAddressRegionGroup),
	)
}

func (updateAddressRegionGroupRepo *approveUpdateAddressRegionGroupRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return updateAddressRegionGroupRepo.approveUpdateAddressRegionGroupTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalUpdateAddressRegionGroup),
	)
}

func (updateAddressRegionGroupRepo *approveUpdateAddressRegionGroupRepository) RunTransaction(
	input *model.InternalUpdateAddressRegionGroup,
) (*model.AddressRegionGroup, error) {
	output, err := updateAddressRegionGroupRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.AddressRegionGroup), err
}
