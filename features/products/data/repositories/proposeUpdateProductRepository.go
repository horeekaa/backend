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

type proposeUpdateProductRepository struct {
	productDataSource                        databaseproductdatasourceinterfaces.ProductDataSource
	createDescriptivePhotoComponent          descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent
	updateDescriptivePhotoComponent          descriptivephotodomainrepositoryinterfaces.UpdateDescriptivePhotoTransactionComponent
	createProductVariantComponent            productvariantdomainrepositoryinterfaces.CreateProductVariantTransactionComponent
	updateProductVariantComponent            productvariantdomainrepositoryinterfaces.UpdateProductVariantTransactionComponent
	createTaggingComponent                   taggingdomainrepositoryinterfaces.CreateTaggingTransactionComponent
	updateTaggingComponent                   taggingdomainrepositoryinterfaces.UpdateTaggingTransactionComponent
	proposeUpdateProductTransactionComponent productdomainrepositoryinterfaces.ProposeUpdateProductTransactionComponent
	mongoDBTransaction                       mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewProposeUpdateProductRepository(
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
	updateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.UpdateDescriptivePhotoTransactionComponent,
	createProductVariantComponent productvariantdomainrepositoryinterfaces.CreateProductVariantTransactionComponent,
	updateProductVariantComponent productvariantdomainrepositoryinterfaces.UpdateProductVariantTransactionComponent,
	createTaggingComponent taggingdomainrepositoryinterfaces.CreateTaggingTransactionComponent,
	updateTaggingComponent taggingdomainrepositoryinterfaces.UpdateTaggingTransactionComponent,
	proposeUpdateProductRepositoryTransactionComponent productdomainrepositoryinterfaces.ProposeUpdateProductTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (productdomainrepositoryinterfaces.ProposeUpdateProductRepository, error) {
	proposeUpdateProductRepo := &proposeUpdateProductRepository{
		productDataSource,
		createDescriptivePhotoComponent,
		updateDescriptivePhotoComponent,
		createProductVariantComponent,
		updateProductVariantComponent,
		createTaggingComponent,
		updateTaggingComponent,
		proposeUpdateProductRepositoryTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		proposeUpdateProductRepo,
		"ProposeUpdateProductRepository",
	)

	return proposeUpdateProductRepo, nil
}

func (updateOrgRepo *proposeUpdateProductRepository) SetValidation(
	usecaseComponent productdomainrepositoryinterfaces.ProposeUpdateProductUsecaseComponent,
) (bool, error) {
	updateOrgRepo.proposeUpdateProductTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (updateOrgRepo *proposeUpdateProductRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateOrgRepo.proposeUpdateProductTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateProduct),
	)
}

func (updateOrgRepo *proposeUpdateProductRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	productToUpdate := input.(*model.InternalUpdateProduct)
	existingProduct, err := updateOrgRepo.productDataSource.GetMongoDataSource().FindByID(
		productToUpdate.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/proposeUpdateProductRepository",
			err,
		)
	}

	if productToUpdate.Photos != nil {
		savedPhotos := existingProduct.Photos
		for _, descPhotoToUpdate := range productToUpdate.Photos {
			if descPhotoToUpdate.ID != nil {
				if !funk.Contains(
					existingProduct.Photos,
					func(dp *model.DescriptivePhoto) bool {
						return dp.ID == *descPhotoToUpdate.ID
					},
				) {
					continue
				}

				_, err := updateOrgRepo.updateDescriptivePhotoComponent.TransactionBody(
					operationOption,
					descPhotoToUpdate,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/proposeUpdateProductRepository",
						err,
					)
				}
				continue
			}

			photoToCreate := &model.InternalCreateDescriptivePhoto{}
			jsonTemp, _ := json.Marshal(descPhotoToUpdate)
			json.Unmarshal(jsonTemp, photoToCreate)
			if descPhotoToUpdate.Photo != nil {
				photoToCreate.Photo.File = descPhotoToUpdate.Photo.File
			}
			photoToCreate.Category = model.DescriptivePhotoCategoryProduct
			photoToCreate.Object = &model.ObjectIDOnly{
				ID: &existingProduct.ID,
			}

			savedPhoto, err := updateOrgRepo.createDescriptivePhotoComponent.TransactionBody(
				operationOption,
				photoToCreate,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/proposeUpdateProductRepository",
					err,
				)
			}
			savedPhotos = append(savedPhotos, savedPhoto)
		}
		if len(savedPhotos) > len(existingProduct.Photos) {
			jsonTemp, _ := json.Marshal(
				map[string]interface{}{
					"Photos": savedPhotos,
				},
			)
			json.Unmarshal(jsonTemp, productToUpdate)
		}
	}

	if productToUpdate.Variants != nil {
		savedVariants := existingProduct.Variants
		for _, variantToUpdate := range productToUpdate.Variants {
			if variantToUpdate.ID != nil {
				if !funk.Contains(
					existingProduct.Variants,
					func(pv *model.ProductVariant) bool {
						return pv.ID == *variantToUpdate.ID
					},
				) {
					continue
				}

				_, err := updateOrgRepo.updateProductVariantComponent.TransactionBody(
					operationOption,
					variantToUpdate,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/proposeUpdateProductRepository",
						err,
					)
				}
				continue
			}

			variantToCreate := &model.InternalCreateProductVariant{}
			jsonTemp, _ := json.Marshal(variantToUpdate)
			json.Unmarshal(jsonTemp, variantToCreate)
			if funk.Get(variantToUpdate, "Photo.Photo") != nil {
				variantToCreate.Photo.Photo.File = variantToUpdate.Photo.Photo.File
			}
			variantToCreate.Product = &model.ObjectIDOnly{
				ID: &existingProduct.ID,
			}

			savedVariant, err := updateOrgRepo.createProductVariantComponent.TransactionBody(
				operationOption,
				variantToCreate,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/proposeUpdateProductRepository",
					err,
				)
			}
			savedVariants = append(savedVariants, savedVariant)
		}
		if len(savedVariants) > len(existingProduct.Variants) {
			jsonTemp, _ := json.Marshal(
				map[string]interface{}{
					"Variants": savedVariants,
				},
			)
			json.Unmarshal(jsonTemp, productToUpdate)
		}
	}

	if productToUpdate.Taggings != nil {
		savedTaggings := existingProduct.Taggings
		for _, taggingToUpdate := range productToUpdate.Taggings {
			if taggingToUpdate.ID != nil {
				if !funk.Contains(
					existingProduct.Taggings,
					func(pv *model.Tagging) bool {
						return pv.ID == *taggingToUpdate.ID
					},
				) {
					continue
				}

				bulkUpdateTagging := &model.InternalBulkUpdateTagging{}
				jsonTemp, _ := json.Marshal(taggingToUpdate)
				json.Unmarshal(jsonTemp, bulkUpdateTagging)
				jsonTemp, _ = json.Marshal(map[string]interface{}{
					"IDs": []interface{}{taggingToUpdate.ID},
				})
				json.Unmarshal(jsonTemp, bulkUpdateTagging)

				_, err := updateOrgRepo.updateTaggingComponent.TransactionBody(
					operationOption,
					bulkUpdateTagging,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/proposeUpdateProductRepository",
						err,
					)
				}
				continue
			}

			taggingToCreate := &model.InternalCreateTagging{}
			jsonTemp, _ := json.Marshal(taggingToUpdate)
			json.Unmarshal(jsonTemp, taggingToCreate)
			taggingToCreate.Products = []*model.ObjectIDOnly{
				{ID: &existingProduct.ID},
			}
			taggingToCreate.TaggingType = func(tt model.TaggingType) *model.TaggingType {
				return &tt
			}(model.TaggingTypeProduct)

			savedTagging, err := updateOrgRepo.createTaggingComponent.TransactionBody(
				operationOption,
				taggingToCreate,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/proposeUpdateProductRepository",
					err,
				)
			}
			savedTaggings = append(savedTaggings, savedTagging...)
		}
		if len(savedTaggings) > len(existingProduct.Taggings) {
			jsonTemp, _ := json.Marshal(
				map[string]interface{}{
					"Taggings": savedTaggings,
				},
			)
			json.Unmarshal(jsonTemp, productToUpdate)
		}
	}

	return updateOrgRepo.proposeUpdateProductTransactionComponent.TransactionBody(
		operationOption,
		productToUpdate,
	)
}

func (updateOrgRepo *proposeUpdateProductRepository) RunTransaction(
	input *model.InternalUpdateProduct,
) (*model.Product, error) {
	output, err := updateOrgRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.Product), err
}
