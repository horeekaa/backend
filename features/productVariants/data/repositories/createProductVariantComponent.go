package productvariantdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	productvariantdomainrepositoryinterfaces "github.com/horeekaa/backend/features/productVariants/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createProductVariantTransactionComponent struct {
	productVariantDataSource        databaseproductvariantdatasourceinterfaces.ProductVariantDataSource
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent
	generatedObjectID               *primitive.ObjectID
}

func (createProdVariantTrx *createProductVariantTransactionComponent) GenerateNewObjectID() primitive.ObjectID {
	generatedObjectID := createProdVariantTrx.productVariantDataSource.GetMongoDataSource().GenerateObjectID()
	createProdVariantTrx.generatedObjectID = &generatedObjectID
	return *createProdVariantTrx.generatedObjectID
}

func (createProdVariantTrx *createProductVariantTransactionComponent) GetCurrentObjectID() primitive.ObjectID {
	if createProdVariantTrx.generatedObjectID == nil {
		generatedObjectID := createProdVariantTrx.productVariantDataSource.GetMongoDataSource().GenerateObjectID()
		createProdVariantTrx.generatedObjectID = &generatedObjectID
	}
	return *createProdVariantTrx.generatedObjectID
}

func NewCreateProductVariantTransactionComponent(
	productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
) (productvariantdomainrepositoryinterfaces.CreateProductVariantTransactionComponent, error) {
	return &createProductVariantTransactionComponent{
		productVariantDataSource:        productVariantDataSource,
		createDescriptivePhotoComponent: createDescriptivePhotoComponent,
	}, nil
}

func (createProdVariantTrx *createProductVariantTransactionComponent) PreTransaction(
	createProductVariantInput *model.InternalCreateProductVariant,
) (*model.InternalCreateProductVariant, error) {
	return createProductVariantInput, nil
}

func (createProdVariantTrx *createProductVariantTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateProductVariant,
) (*model.ProductVariant, error) {
	variantToCreate := &model.DatabaseCreateProductVariant{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, variantToCreate)
	variantToCreate.ID = createProdVariantTrx.GetCurrentObjectID()

	if input.Photo != nil {
		input.Photo.Category = model.DescriptivePhotoCategoryProductVariant
		input.Photo.Object = &model.ObjectIDOnly{
			ID: &variantToCreate.ID,
		}
		descriptivePhoto, err := createProdVariantTrx.createDescriptivePhotoComponent.TransactionBody(
			&mongodbcoretypes.OperationOptions{},
			input.Photo,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/createProductVariant",
				err,
			)
		}

		variantToCreate.Photo = &model.ObjectIDOnly{
			ID: &descriptivePhoto.ID,
		}
	}

	createdVariant, err := createProdVariantTrx.productVariantDataSource.GetMongoDataSource().Create(
		variantToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createProductVariant",
			err,
		)
	}
	createProdVariantTrx.generatedObjectID = nil

	return createdVariant, nil
}
