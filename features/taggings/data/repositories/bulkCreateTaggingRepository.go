package taggingdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type bulkCreateTaggingRepository struct {
	bulkCreateTaggingTransactionComponent taggingdomainrepositoryinterfaces.BulkCreateTaggingTransactionComponent
	mongoDBTransaction                    mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewBulkCreateTaggingRepository(
	bulkCreateTaggingRepositoryTransactionComponent taggingdomainrepositoryinterfaces.BulkCreateTaggingTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (taggingdomainrepositoryinterfaces.BulkCreateTaggingRepository, error) {
	bulkCreateTaggingRepo := &bulkCreateTaggingRepository{
		bulkCreateTaggingRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		bulkCreateTaggingRepo,
		"BulkCreateTaggingRepository",
	)

	return bulkCreateTaggingRepo, nil
}

func (bulkCreateTaggingRepo *bulkCreateTaggingRepository) SetValidation(
	usecaseComponent taggingdomainrepositoryinterfaces.BulkCreateTaggingUsecaseComponent,
) (bool, error) {
	bulkCreateTaggingRepo.bulkCreateTaggingTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (bulkCreateTaggingRepo *bulkCreateTaggingRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return bulkCreateTaggingRepo.bulkCreateTaggingTransactionComponent.PreTransaction(
		input.(*model.InternalCreateTagging),
	)
}

func (bulkCreateTaggingRepo *bulkCreateTaggingRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	return bulkCreateTaggingRepo.bulkCreateTaggingTransactionComponent.TransactionBody(
		operationOption,
		input.(*model.InternalCreateTagging),
	)
}

func (bulkCreateTaggingRepo *bulkCreateTaggingRepository) RunTransaction(
	input *model.InternalCreateTagging,
) ([]*model.Tagging, error) {
	output, err := bulkCreateTaggingRepo.mongoDBTransaction.RunTransaction(input)
	return (output).([]*model.Tagging), err
}
