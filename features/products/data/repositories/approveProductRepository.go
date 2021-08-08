package productdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type approveUpdateProductRepository struct {
	approveUpdateProductTransactionComponent productdomainrepositoryinterfaces.ApproveUpdateProductTransactionComponent
	productDataSource                        databaseproductdatasourceinterfaces.ProductDataSource
	bulkApproveUpdateTaggingComponent        taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingTransactionComponent
	mongoDBTransaction                       mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewApproveUpdateProductRepository(
	approveUpdateProductRepositoryTransactionComponent productdomainrepositoryinterfaces.ApproveUpdateProductTransactionComponent,
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	bulkApproveUpdateTaggingComponent taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (productdomainrepositoryinterfaces.ApproveUpdateProductRepository, error) {
	approveUpdateProductRepo := &approveUpdateProductRepository{
		approveUpdateProductRepositoryTransactionComponent,
		productDataSource,
		bulkApproveUpdateTaggingComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		approveUpdateProductRepo,
		"ApproveUpdateProductRepository",
	)

	return approveUpdateProductRepo, nil
}

func (updateProdRepo *approveUpdateProductRepository) SetValidation(
	usecaseComponent productdomainrepositoryinterfaces.ApproveUpdateProductUsecaseComponent,
) (bool, error) {
	updateProdRepo.approveUpdateProductTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (updateProdRepo *approveUpdateProductRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateProdRepo.approveUpdateProductTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateProduct),
	)
}

func (updateProdRepo *approveUpdateProductRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	productToApprove := input.(*model.InternalUpdateProduct)
	existingProduct, err := updateProdRepo.productDataSource.GetMongoDataSource().FindByID(
		productToApprove.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/approveUpdateProductRepository",
			err,
		)
	}

	if existingProduct.ProposedChanges.ProposalStatus == model.EntityProposalStatusProposed {
		bulkUpdateTagging := &model.InternalBulkUpdateTagging{}
		jsonTemp, _ := json.Marshal(map[string]interface{}{
			"IDs": funk.Map(
				existingProduct.ProposedChanges.Taggings,
				func(_, tagging *model.Tag) interface{} {
					return tagging.ID
				},
			),
		})
		json.Unmarshal(jsonTemp, bulkUpdateTagging)

		bulkUpdateTagging.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: productToApprove.RecentApprovingAccount.ID,
		}
		bulkUpdateTagging.ProposalStatus = productToApprove.ProposalStatus

		updateProdRepo.bulkApproveUpdateTaggingComponent.TransactionBody(
			operationOption,
			bulkUpdateTagging,
		)
	}

	return updateProdRepo.approveUpdateProductTransactionComponent.TransactionBody(
		operationOption,
		productToApprove,
	)
}

func (updateProdRepo *approveUpdateProductRepository) RunTransaction(
	input *model.InternalUpdateProduct,
) (*model.Product, error) {
	output, err := updateProdRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.Product), err
}
