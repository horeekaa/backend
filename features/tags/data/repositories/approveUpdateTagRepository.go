package tagdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	tagdomainrepositoryinterfaces "github.com/horeekaa/backend/features/tags/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateTagRepository struct {
	approveUpdateTagTransactionComponent tagdomainrepositoryinterfaces.ApproveUpdateTagTransactionComponent
	mongoDBTransaction                   mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewApproveUpdateTagRepository(
	approveUpdateTagRepositoryTransactionComponent tagdomainrepositoryinterfaces.ApproveUpdateTagTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (tagdomainrepositoryinterfaces.ApproveUpdateTagRepository, error) {
	approveUpdatetagRepo := &approveUpdateTagRepository{
		approveUpdateTagRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		approveUpdatetagRepo,
		"ApproveUpdateTagRepository",
	)

	return approveUpdatetagRepo, nil
}

func (updateTagRepo *approveUpdateTagRepository) SetValidation(
	usecaseComponent tagdomainrepositoryinterfaces.ApproveUpdateTagUsecaseComponent,
) (bool, error) {
	updateTagRepo.approveUpdateTagTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (updateTagRepo *approveUpdateTagRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateTagRepo.approveUpdateTagTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateTag),
	)
}

func (updateTagRepo *approveUpdateTagRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return updateTagRepo.approveUpdateTagTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalUpdateTag),
	)
}

func (updateTagRepo *approveUpdateTagRepository) RunTransaction(
	input *model.InternalUpdateTag,
) (*model.Tag, error) {
	output, err := updateTagRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.Tag), err
}
