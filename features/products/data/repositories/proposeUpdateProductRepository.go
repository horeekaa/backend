package productdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type proposeUpdateProductRepository struct {
	productDataSource                        databaseproductdatasourceinterfaces.ProductDataSource
	createDescriptivePhotoComponent          descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent
	updateDescriptivePhotoComponent          descriptivephotodomainrepositoryinterfaces.UpdateDescriptivePhotoTransactionComponent
	proposeUpdateProductTransactionComponent productdomainrepositoryinterfaces.ProposeUpdateProductTransactionComponent
	mongoDBTransaction                       mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewProposeUpdateProductRepository(
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
	updateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.UpdateDescriptivePhotoTransactionComponent,
	proposeUpdateProductRepositoryTransactionComponent productdomainrepositoryinterfaces.ProposeUpdateProductTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (productdomainrepositoryinterfaces.ProposeUpdateProductRepository, error) {
	proposeUpdateProductRepo := &proposeUpdateProductRepository{
		productDataSource,
		createDescriptivePhotoComponent,
		updateDescriptivePhotoComponent,
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
	if productToUpdate.Photos != nil {
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
			photoToCreate.Photo.File = descPhotoToUpdate.Photo.File
			photoToCreate.Category = model.DescriptivePhotoCategoryProduct

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
