package taggingdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type bulkApproveUpdateTaggingRepository struct {
	bulkApproveUpdateTaggingTransactionComponent taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingTransactionComponent
	mongoDBTransaction                           mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewBulkApproveUpdateTaggingRepository(
	bulkApproveUpdateTaggingRepositoryTransactionComponent taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingRepository, error) {
	bulkApproveUpdatetaggingRepo := &bulkApproveUpdateTaggingRepository{
		bulkApproveUpdateTaggingRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		bulkApproveUpdatetaggingRepo,
		"BulkApproveUpdateTaggingRepository",
	)

	return bulkApproveUpdatetaggingRepo, nil
}

func (bulkUpdateTaggingRepo *bulkApproveUpdateTaggingRepository) SetValidation(
	usecaseComponent taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingUsecaseComponent,
) (bool, error) {
	bulkUpdateTaggingRepo.bulkApproveUpdateTaggingTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (bulkApproveUpdateTaggingRepo *bulkApproveUpdateTaggingRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return bulkApproveUpdateTaggingRepo.bulkApproveUpdateTaggingTransactionComponent.PreTransaction(
		input.(*model.InternalBulkUpdateTagging),
	)
}

func (bulkApproveUpdateTaggingRepo *bulkApproveUpdateTaggingRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return bulkApproveUpdateTaggingRepo.bulkApproveUpdateTaggingTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalBulkUpdateTagging),
	)
}

func (bulkApproveUpdateTaggingRepo *bulkApproveUpdateTaggingRepository) RunTransaction(
	input *model.InternalBulkUpdateTagging,
) ([]*model.Tagging, error) {
	output, err := bulkApproveUpdateTaggingRepo.mongoDBTransaction.RunTransaction(input)
	return (output).([]*model.Tagging), err
}
