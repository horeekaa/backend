package productvariantdomainrepositories

import (
	"encoding/json"
	"time"

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
	pathIdentity                    string
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
		pathIdentity:                    "CreateProductVariantComponent",
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
	productVariantToCreate := &model.DatabaseCreateProductVariant{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, productVariantToCreate)

	generatedObjectID := createProdVariantTrx.GetCurrentObjectID()
	if input.Photo != nil {
		input.Photo.Category = model.DescriptivePhotoCategoryProductVariant
		input.Photo.Object = &model.ObjectIDOnly{
			ID: &generatedObjectID,
		}
		input.Photo.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
			return &s
		}(*input.ProposalStatus)
		input.Photo.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
			return &m
		}(*input.SubmittingAccount)
		descriptivePhoto, err := createProdVariantTrx.createDescriptivePhotoComponent.TransactionBody(
			session,
			input.Photo,
		)
		if err != nil {
			return nil, err
		}

		productVariantToCreate.Photo = &model.ObjectIDOnly{
			ID: &descriptivePhoto.ID,
		}
	}

	newDocumentJson, _ := json.Marshal(*productVariantToCreate)
	loggingOutput, err := createProdVariantTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "ProductVariant",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: productVariantToCreate.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *productVariantToCreate.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createProdVariantTrx.pathIdentity,
			err,
		)
	}

	productVariantToCreate.ID = generatedObjectID
	productVariantToCreate.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *productVariantToCreate.ProposalStatus == model.EntityProposalStatusApproved {
		productVariantToCreate.RecentApprovingAccount = &model.ObjectIDOnly{ID: productVariantToCreate.SubmittingAccount.ID}
	}

	currentTime := time.Now().UTC()
	productVariantToCreate.CreatedAt = &currentTime
	productVariantToCreate.UpdatedAt = &currentTime

	defaultIsActive := true
	if productVariantToCreate.IsActive == nil {
		productVariantToCreate.IsActive = &defaultIsActive
	}

	defaultProposalStatus := model.EntityProposalStatusProposed
	if productVariantToCreate.ProposalStatus == nil {
		productVariantToCreate.ProposalStatus = &defaultProposalStatus
	}

	jsonTemp, _ = json.Marshal(productVariantToCreate)
	json.Unmarshal(jsonTemp, &productVariantToCreate.ProposedChanges)

	createdVariant, err := createProdVariantTrx.productVariantDataSource.GetMongoDataSource().Create(
		productVariantToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createProdVariantTrx.pathIdentity,
			err,
		)
	}
	createProdVariantTrx.generatedObjectID = nil

	return createdVariant, nil
}
