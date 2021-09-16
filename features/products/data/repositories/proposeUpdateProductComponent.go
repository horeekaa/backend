package productdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseproductdatasourceinterfaces "github.com/horeekaa/backend/features/products/data/dataSources/databases/interfaces/sources"
	productdomainrepositoryinterfaces "github.com/horeekaa/backend/features/products/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type proposeUpdateProductTransactionComponent struct {
	productDataSource                    databaseproductdatasourceinterfaces.ProductDataSource
	loggingDataSource                    databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility                  coreutilityinterfaces.MapProcessorUtility
	proposeUpdateProductUsecaseComponent productdomainrepositoryinterfaces.ProposeUpdateProductUsecaseComponent
}

func NewProposeUpdateProductTransactionComponent(
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (productdomainrepositoryinterfaces.ProposeUpdateProductTransactionComponent, error) {
	return &proposeUpdateProductTransactionComponent{
		productDataSource:   productDataSource,
		loggingDataSource:   loggingDataSource,
		mapProcessorUtility: mapProcessorUtility,
	}, nil
}

func (updateProdTrx *proposeUpdateProductTransactionComponent) SetValidation(
	usecaseComponent productdomainrepositoryinterfaces.ProposeUpdateProductUsecaseComponent,
) (bool, error) {
	updateProdTrx.proposeUpdateProductUsecaseComponent = usecaseComponent
	return true, nil
}

func (updateProdTrx *proposeUpdateProductTransactionComponent) PreTransaction(
	input *model.InternalUpdateProduct,
) (*model.InternalUpdateProduct, error) {
	if updateProdTrx.proposeUpdateProductUsecaseComponent == nil {
		return input, nil
	}
	return updateProdTrx.proposeUpdateProductUsecaseComponent.Validation(input)
}

func (updateProdTrx *proposeUpdateProductTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateProduct *model.InternalUpdateProduct,
) (*model.Product, error) {
	existingProduct, err := updateProdTrx.productDataSource.GetMongoDataSource().FindByID(
		updateProduct.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateProduct",
			err,
		)
	}

	newDocumentJson, _ := json.Marshal(*updateProduct)
	oldDocumentJson, _ := json.Marshal(*existingProduct)
	loggingOutput, err := updateProdTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "Product",
			Document: &model.ObjectIDOnly{
				ID: &existingProduct.ID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			OldDocumentJSON: func(s string) *string { return &s }(string(oldDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updateProduct.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updateProduct.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateProduct",
			err,
		)
	}
	updateProduct.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	fieldsToUpdateProduct := &model.DatabaseUpdateProduct{
		ID: updateProduct.ID,
	}
	jsonExisting, _ := json.Marshal(existingProduct)
	json.Unmarshal(jsonExisting, &fieldsToUpdateProduct.ProposedChanges)

	var updateProductMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateProduct)
	json.Unmarshal(jsonUpdate, &updateProductMap)

	updateProdTrx.mapProcessorUtility.RemoveNil(updateProductMap)

	jsonUpdate, _ = json.Marshal(updateProductMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateProduct.ProposedChanges)

	if updateProduct.ProposalStatus != nil {
		fieldsToUpdateProduct.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: updateProduct.SubmittingAccount.ID,
		}
		if *updateProduct.ProposalStatus == model.EntityProposalStatusApproved {
			json.Unmarshal(jsonUpdate, fieldsToUpdateProduct)
		}
	}

	updatedProduct, err := updateProdTrx.productDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateProduct.ID,
		},
		fieldsToUpdateProduct,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateProduct",
			err,
		)
	}

	return updatedProduct, nil
}
