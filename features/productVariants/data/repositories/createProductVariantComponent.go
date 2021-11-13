package productvariantdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	productvariantdomainrepositoryinterfaces "github.com/horeekaa/backend/features/productVariants/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createProductVariantTransactionComponent struct {
	productVariantDataSource        databaseproductvariantdatasourceinterfaces.ProductVariantDataSource
	loggingDataSource               databaseloggingdatasourceinterfaces.LoggingDataSource
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
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
) (productvariantdomainrepositoryinterfaces.CreateProductVariantTransactionComponent, error) {
	return &createProductVariantTransactionComponent{
		productVariantDataSource:        productVariantDataSource,
		loggingDataSource:               loggingDataSource,
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
	if input.Photo != nil {
		input.Photo.Category = model.DescriptivePhotoCategoryProductVariant
		input.Photo.Object = &model.ObjectIDOnly{
			ID: input.ID,
		}
		input.Photo.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
			return &s
		}(*input.ProposalStatus)
		input.Photo.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
			return &m
		}(*input.SubmittingAccount)
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

		input.Photo = &model.InternalCreateDescriptivePhoto{
			ID: &descriptivePhoto.ID,
		}
	}

	newDocumentJson, _ := json.Marshal(*input)
	generatedObjectID := createProdVariantTrx.GetCurrentObjectID()
	loggingOutput, err := createProdVariantTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "ProductVariant",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: input.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *input.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createProductVariant",
			err,
		)
	}

	input.ID = &generatedObjectID
	input.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *input.ProposalStatus == model.EntityProposalStatusApproved {
		input.RecentApprovingAccount = &model.ObjectIDOnly{ID: input.SubmittingAccount.ID}
	}

	variantToCreate := &model.DatabaseCreateProductVariant{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, variantToCreate)
	json.Unmarshal(jsonTemp, &variantToCreate.ProposedChanges)

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
