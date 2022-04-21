package productdomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createProductTransactionComponent struct {
	productDataSource             databaseproductdatasourceinterfaces.ProductDataSource
	loggingDataSource             databaseloggingdatasourceinterfaces.LoggingDataSource
	createProductUsecaseComponent productdomainrepositoryinterfaces.CreateProductUsecaseComponent
	generatedObjectID             *primitive.ObjectID
	pathIdentity                  string
}

func NewCreateProductTransactionComponent(
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
) (productdomainrepositoryinterfaces.CreateProductTransactionComponent, error) {
	return &createProductTransactionComponent{
		productDataSource: productDataSource,
		loggingDataSource: loggingDataSource,
		pathIdentity:      "CreateProductComponent",
	}, nil
}

func (createProductTrx *createProductTransactionComponent) GenerateNewObjectID() primitive.ObjectID {
	generatedObjectID := createProductTrx.productDataSource.GetMongoDataSource().GenerateObjectID()
	createProductTrx.generatedObjectID = &generatedObjectID
	return *createProductTrx.generatedObjectID
}

func (createProductTrx *createProductTransactionComponent) GetCurrentObjectID() primitive.ObjectID {
	if createProductTrx.generatedObjectID == nil {
		generatedObjectID := createProductTrx.productDataSource.GetMongoDataSource().GenerateObjectID()
		createProductTrx.generatedObjectID = &generatedObjectID
	}
	return *createProductTrx.generatedObjectID
}

func (createProductTrx *createProductTransactionComponent) SetValidation(
	usecaseComponent productdomainrepositoryinterfaces.CreateProductUsecaseComponent,
) (bool, error) {
	createProductTrx.createProductUsecaseComponent = usecaseComponent
	return true, nil
}

func (createProductTrx *createProductTransactionComponent) PreTransaction(
	input *model.InternalCreateProduct,
) (*model.InternalCreateProduct, error) {
	if createProductTrx.createProductUsecaseComponent == nil {
		return input, nil
	}
	return createProductTrx.createProductUsecaseComponent.Validation(input)
}

func (createProductTrx *createProductTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalCreateProduct,
) (*model.Product, error) {
	productToCreate := &model.DatabaseCreateProduct{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, productToCreate)

	newDocumentJson, _ := json.Marshal(*productToCreate)
	generatedObjectID := createProductTrx.GetCurrentObjectID()
	loggingOutput, err := createProductTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Product",
			Document: &model.ObjectIDOnly{
				ID: &generatedObjectID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: productToCreate.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityCreate,
			ProposalStatus: *productToCreate.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createProductTrx.pathIdentity,
			err,
		)
	}

	productToCreate.ID = generatedObjectID
	productToCreate.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}
	if *productToCreate.ProposalStatus == model.EntityProposalStatusApproved {
		productToCreate.RecentApprovingAccount = &model.ObjectIDOnly{ID: productToCreate.SubmittingAccount.ID}
	}

	currentTime := time.Now().UTC()
	productToCreate.CreatedAt = &currentTime
	productToCreate.UpdatedAt = &currentTime

	defaultProposalStatus := model.EntityProposalStatusProposed
	if productToCreate.ProposalStatus == nil {
		productToCreate.ProposalStatus = &defaultProposalStatus
	}

	defaultIsActive := true
	if productToCreate.IsActive == nil {
		productToCreate.IsActive = &defaultIsActive
	}
	if productToCreate.Photos == nil {
		productToCreate.Photos = []*model.ObjectIDOnly{}
	}
	if productToCreate.Variants == nil {
		productToCreate.Variants = []*model.ObjectIDOnly{}
	}
	if productToCreate.Taggings == nil {
		productToCreate.Taggings = []*model.ObjectIDOnly{}
	}

	jsonTemp, _ = json.Marshal(productToCreate)
	json.Unmarshal(jsonTemp, &productToCreate.ProposedChanges)

	newProduct, err := createProductTrx.productDataSource.GetMongoDataSource().Create(
		productToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			createProductTrx.pathIdentity,
			err,
		)
	}
	createProductTrx.generatedObjectID = nil

	return newProduct, nil
}
