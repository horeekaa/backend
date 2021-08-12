package taggingdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type bulkProposeUpdateTaggingRepository struct {
	bulkProposeUpdateTaggingTransactionComponent taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingTransactionComponent
	mongoDBTransaction                           mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewBulkProposeUpdateTaggingRepository(
	bulkProposeUpdateTaggingRepositoryTransactionComponent taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingRepository, error) {
	bulkProposeUpdateTaggingRepo := &bulkProposeUpdateTaggingRepository{
		bulkProposeUpdateTaggingRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		bulkProposeUpdateTaggingRepo,
		"BulkProposeUpdateTaggingRepository",
	)

	return bulkProposeUpdateTaggingRepo, nil
}

func (bulkUpdateTaggingRepo *bulkProposeUpdateTaggingRepository) SetValidation(
	usecaseComponent taggingdomainrepositoryinterfaces.BulkProposeUpdateTaggingUsecaseComponent,
) (bool, error) {
	bulkUpdateTaggingRepo.bulkProposeUpdateTaggingTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (bulkUpdateTaggingRepo *bulkProposeUpdateTaggingRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return bulkUpdateTaggingRepo.bulkProposeUpdateTaggingTransactionComponent.PreTransaction(
		input.(*model.InternalBulkUpdateTagging),
	)
}

func (bulkUpdateTaggingRepo *bulkProposeUpdateTaggingRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return bulkUpdateTaggingRepo.bulkProposeUpdateTaggingTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalBulkUpdateTagging),
	)
}

func (bulkUpdateTaggingRepo *bulkProposeUpdateTaggingRepository) RunTransaction(
	input *model.InternalBulkUpdateTagging,
) ([]*model.Tagging, error) {
	output, err := bulkUpdateTaggingRepo.mongoDBTransaction.RunTransaction(input)
	return (output).([]*model.Tagging), err
}
