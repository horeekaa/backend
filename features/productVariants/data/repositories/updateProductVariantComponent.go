package productvariantdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	productvariantdomainrepositoryinterfaces "github.com/horeekaa/backend/features/productVariants/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type updateProductVariantTransactionComponent struct {
	productVariantDataSource        databaseproductvariantdatasourceinterfaces.ProductVariantDataSource
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent
	updateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.UpdateDescriptivePhotoTransactionComponent
}

func NewUpdateProductVariantTransactionComponent(
	productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
	updateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.UpdateDescriptivePhotoTransactionComponent,
) (productvariantdomainrepositoryinterfaces.UpdateProductVariantTransactionComponent, error) {
	return &updateProductVariantTransactionComponent{
		productVariantDataSource:        productVariantDataSource,
		createDescriptivePhotoComponent: createDescriptivePhotoComponent,
		updateDescriptivePhotoComponent: updateDescriptivePhotoComponent,
	}, nil
}

func (updateDescPhotoTrx *updateProductVariantTransactionComponent) PreTransaction(
	updateProductVariantInput *model.InternalUpdateProductVariant,
) (*model.InternalUpdateProductVariant, error) {
	return updateProductVariantInput, nil
}

func (updateDescPhotoTrx *updateProductVariantTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateProductVariant,
) (*model.ProductVariant, error) {
	_, err := updateDescPhotoTrx.productVariantDataSource.GetMongoDataSource().FindByID(
		*input.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateProductVariant",
			err,
		)
	}

	productVariantToUpdate := &model.DatabaseUpdateProductVariant{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, productVariantToUpdate)

	if input.Photo != nil {
		if input.Photo.ID != nil {
			_, err := updateDescPhotoTrx.updateDescriptivePhotoComponent.TransactionBody(
				session,
				input.Photo,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/updateProductVariant",
					err,
				)
			}
		} else {
			photoToCreate := &model.InternalCreateDescriptivePhoto{}
			jsonTemp, _ := json.Marshal(input)
			json.Unmarshal(jsonTemp, photoToCreate)
			photoToCreate.Category = model.DescriptivePhotoCategoryProductVariant
			photoToCreate.Photo.File = input.Photo.Photo.File
			createdDescriptivePhoto, err := updateDescPhotoTrx.createDescriptivePhotoComponent.TransactionBody(
				session,
				photoToCreate,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/updateProductVariant",
					err,
				)
			}
			productVariantToUpdate.Photo = &model.ObjectIDOnly{
				ID: &createdDescriptivePhoto.ID,
			}
		}
	}

	updatedDescPhoto, err := updateDescPhotoTrx.productVariantDataSource.GetMongoDataSource().Update(
		productVariantToUpdate.ID,
		productVariantToUpdate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateProductVariant",
			err,
		)
	}

	return updatedDescPhoto, nil
}
