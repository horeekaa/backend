package productdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createProductRepository struct {
	createProductTransactionComponent productdomainrepositoryinterfaces.CreateProductTransactionComponent
	createDescriptivePhotoComponent   descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent
	mongoDBTransaction                mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreateProductRepository(
	createProductRepositoryTransactionComponent productdomainrepositoryinterfaces.CreateProductTransactionComponent,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (productdomainrepositoryinterfaces.CreateProductRepository, error) {
	createProductRepo := &createProductRepository{
		createProductRepositoryTransactionComponent,
		createDescriptivePhotoComponent,
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
	if productToCreate.Photos != nil {
		savedPhotos := []*model.InternalCreateDescriptivePhoto{}
		for _, photo := range productToCreate.Photos {
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
