package productdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	productvariantdomainrepositoryinterfaces "github.com/horeekaa/backend/features/productVariants/domain/repositories"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	taggingdomainrepositoryinterfaces "github.com/horeekaa/backend/features/taggings/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createProductRepository struct {
	createProductTransactionComponent productdomainrepositoryinterfaces.CreateProductTransactionComponent
	createProductVariantComponent     productvariantdomainrepositoryinterfaces.CreateProductVariantTransactionComponent
	createDescriptivePhotoComponent   descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent
	createTaggingComponent            taggingdomainrepositoryinterfaces.CreateTaggingTransactionComponent
	mongoDBTransaction                mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreateProductRepository(
	createProductRepositoryTransactionComponent productdomainrepositoryinterfaces.CreateProductTransactionComponent,
	createProductVariantComponent productvariantdomainrepositoryinterfaces.CreateProductVariantTransactionComponent,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
	createTaggingComponent taggingdomainrepositoryinterfaces.CreateTaggingTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (productdomainrepositoryinterfaces.CreateProductRepository, error) {
	createProductRepo := &createProductRepository{
		createProductRepositoryTransactionComponent,
		createProductVariantComponent,
		createDescriptivePhotoComponent,
		createTaggingComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		createProductRepo,
		"CreateProductRepository",
	)

	return createProductRepo, nil
}

func (createProdRepo *createProductRepository) SetValidation(
	usecaseComponent productdomainrepositoryinterfaces.CreateProductUsecaseComponent,
) (bool, error) {
	createProdRepo.createProductTransactionComponent.SetValidation(usecaseComponent)
	return true, nil
}

func (createProdRepo *createProductRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return createProdRepo.createProductTransactionComponent.PreTransaction(
		input.(*model.InternalCreateProduct),
	)
}

func (createProdRepo *createProductRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	productToCreate := input.(*model.InternalCreateProduct)
	generatedObjectID := createProdRepo.createProductTransactionComponent.GenerateNewObjectID()
	if productToCreate.Photos != nil {
		savedPhotos := []*model.InternalCreateDescriptivePhoto{}
		for _, photo := range productToCreate.Photos {
			photo.Category = model.DescriptivePhotoCategoryProduct
			photo.Object = &model.ObjectIDOnly{
				ID: &generatedObjectID,
			}
			createdPhotoOutput, err := createProdRepo.createDescriptivePhotoComponent.TransactionBody(
				operationOption,
				photo,
			)
			if err != nil {
				return nil, err
			}
			savedPhoto := &model.InternalCreateDescriptivePhoto{}
			jsonTemp, _ := json.Marshal(createdPhotoOutput)
			json.Unmarshal(jsonTemp, savedPhoto)
			savedPhotos = append(savedPhotos, savedPhoto)
		}
		productToCreate.Photos = savedPhotos
	}

	if productToCreate.Variants != nil {
		savedVariants := []*model.InternalCreateProductVariant{}
		for _, variant := range productToCreate.Variants {
			variant.Product = &model.ObjectIDOnly{
				ID: &generatedObjectID,
			}
			createdVariantOutput, err := createProdRepo.createProductVariantComponent.TransactionBody(
				operationOption,
				variant,
			)
			if err != nil {
				return nil, err
			}
			savedVariant := &model.InternalCreateProductVariant{}
			jsonTemp, _ := json.Marshal(createdVariantOutput)
			json.Unmarshal(jsonTemp, savedVariant)
			savedVariants = append(savedVariants, savedVariant)
		}
		productToCreate.Variants = savedVariants
	}

	if productToCreate.Taggings != nil {
		savedTaggings := []*model.InternalCreateTagging{}
		for _, tagging := range productToCreate.Taggings {
			tagging.Products = []*model.ObjectIDOnly{
				{ID: &generatedObjectID},
			}
			createdTaggingOutput, err := createProdRepo.createTaggingComponent.TransactionBody(
				operationOption,
				tagging,
			)
			if err != nil {
				return nil, err
			}

			savedTagging := &model.InternalCreateTagging{}
			jsonTemp, _ := json.Marshal(createdTaggingOutput[0])
			json.Unmarshal(jsonTemp, savedTagging)
			savedTaggings = append(savedTaggings, savedTagging)
		}
		productToCreate.Taggings = savedTaggings
	}

	return createProdRepo.createProductTransactionComponent.TransactionBody(
		operationOption,
		productToCreate,
	)
}

func (createProdRepo *createProductRepository) RunTransaction(
	input *model.InternalCreateProduct,
) (*model.Product, error) {
	output, err := createProdRepo.mongoDBTransaction.RunTransaction(input)
	return (output).(*model.Product), err
}
