package productdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	productvariantdomainrepositoryinterfaces "github.com/horeekaa/backend/features/productVariants/domain/repositories"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type approveUpdateProductRepository struct {
	approveUpdateProductTransactionComponent productdomainrepositoryinterfaces.ApproveUpdateProductTransactionComponent
	productDataSource                        databaseproductdatasourceinterfaces.ProductDataSource
	approveUpdateProductVariantComponent     productvariantdomainrepositoryinterfaces.ApproveUpdateProductVariantTransactionComponent
	approveDescriptivePhotoComponent         descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent
	bulkApproveUpdateTaggingComponent        taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingTransactionComponent
	mongoDBTransaction                       mongodbcoretransactioninterfaces.MongoRepoTransaction
	pathIdentity                             string
}

func NewApproveUpdateProductRepository(
	approveUpdateProductRepositoryTransactionComponent productdomainrepositoryinterfaces.ApproveUpdateProductTransactionComponent,
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	approveUpdateProductVariantComponent productvariantdomainrepositoryinterfaces.ApproveUpdateProductVariantTransactionComponent,
	approveDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent,
	bulkApproveUpdateTaggingComponent taggingdomainrepositoryinterfaces.BulkApproveUpdateTaggingTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (productdomainrepositoryinterfaces.ApproveUpdateProductRepository, error) {
	approveUpdateProductRepo := &approveUpdateProductRepository{
		approveUpdateProductRepositoryTransactionComponent,
		productDataSource,
		approveUpdateProductVariantComponent,
		approveDescriptivePhotoComponent,
		bulkApproveUpdateTaggingComponent,
		mongoDBTransaction,
		"ApproveUpdateProductRepository",
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
			updateProdRepo.pathIdentity,
			err,
		)
	}

	if existingProduct.ProposedChanges.ProposalStatus == model.EntityProposalStatusProposed {
		if existingProduct.ProposedChanges.Photos != nil {
			for _, photo := range existingProduct.ProposedChanges.Photos {
				updateDescriptivePhoto := &model.InternalUpdateDescriptivePhoto{
					ID: &photo.ID,
				}
				updateDescriptivePhoto.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*productToApprove.RecentApprovingAccount)
				updateDescriptivePhoto.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*productToApprove.ProposalStatus)

				_, err := updateProdRepo.approveDescriptivePhotoComponent.TransactionBody(
					operationOption,
					updateDescriptivePhoto,
				)
				if err != nil {
					return nil, err
				}
			}
		}

		if existingProduct.ProposedChanges.Taggings != nil {
			bulkUpdateTagging := &model.InternalBulkUpdateTagging{}
			jsonTemp, _ := json.Marshal(map[string]interface{}{
				"IDs": funk.Map(
					existingProduct.ProposedChanges.Taggings,
					func(_, tagging *model.Tagging) interface{} {
						return tagging.ID
					},
				),
			})
			json.Unmarshal(jsonTemp, bulkUpdateTagging)

			bulkUpdateTagging.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*productToApprove.RecentApprovingAccount)
			bulkUpdateTagging.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*productToApprove.ProposalStatus)

			_, err := updateProdRepo.bulkApproveUpdateTaggingComponent.TransactionBody(
				operationOption,
				bulkUpdateTagging,
			)
			if err != nil {
				return nil, err
			}
		}

		if existingProduct.ProposedChanges.Variants != nil {
			for _, variant := range existingProduct.ProposedChanges.Variants {
				updateVariant := &model.InternalUpdateProductVariant{
					ID: &variant.ID,
				}
				updateVariant.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*productToApprove.RecentApprovingAccount)
				updateVariant.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*productToApprove.ProposalStatus)

				_, err := updateProdRepo.approveUpdateProductVariantComponent.TransactionBody(
					operationOption,
					updateVariant,
				)
				if err != nil {
					return nil, err
				}
			}
		}
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
	if err != nil {
		return nil, err
	}
	return (output).(*model.Product), err
}
