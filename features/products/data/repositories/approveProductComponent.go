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

type approveUpdateProductTransactionComponent struct {
	productDataSource                    databaseproductdatasourceinterfaces.ProductDataSource
	loggingDataSource                    databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility                  coreutilityinterfaces.MapProcessorUtility
	approveUpdateProductUsecaseComponent productdomainrepositoryinterfaces.ApproveUpdateProductUsecaseComponent
}

func NewApproveUpdateProductTransactionComponent(
	productDataSource databaseproductdatasourceinterfaces.ProductDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (productdomainrepositoryinterfaces.ApproveUpdateProductTransactionComponent, error) {
	return &approveUpdateProductTransactionComponent{
		productDataSource:   productDataSource,
		loggingDataSource:   loggingDataSource,
		mapProcessorUtility: mapProcessorUtility,
	}, nil
}

func (approveProdTrx *approveUpdateProductTransactionComponent) SetValidation(
	usecaseComponent productdomainrepositoryinterfaces.ApproveUpdateProductUsecaseComponent,
) (bool, error) {
	approveProdTrx.approveUpdateProductUsecaseComponent = usecaseComponent
	return true, nil
}

func (approveProdTrx *approveUpdateProductTransactionComponent) PreTransaction(
	input *model.InternalUpdateProduct,
) (*model.InternalUpdateProduct, error) {
	if approveProdTrx.approveUpdateProductUsecaseComponent == nil {
		return input, nil
	}
	return approveProdTrx.approveUpdateProductUsecaseComponent.Validation(input)
}

func (approveProdTrx *approveUpdateProductTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateProduct *model.InternalUpdateProduct,
) (*model.Product, error) {
	existingProduct, err := approveProdTrx.productDataSource.GetMongoDataSource().FindByID(
		updateProduct.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateProduct",
			err,
		)
	}

	previousLog, err := approveProdTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingProduct.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateProduct",
			err,
		)
	}

	logToCreate := &model.CreateLogging{
		Collection: previousLog.Collection,
		CreatedByAccount: &model.ObjectIDOnly{
			ID: updateProduct.RecentApprovingAccount.ID,
		},
		Activity:       previousLog.Activity,
		ProposalStatus: *updateProduct.ProposalStatus,
	}
	jsonTemp, _ := json.Marshal(
		map[string]interface{}{
			"FieldChanges": previousLog.FieldChanges,
		},
	)
	json.Unmarshal(jsonTemp, logToCreate)

	createdLog, err := approveProdTrx.loggingDataSource.GetMongoDataSource().Create(
		logToCreate,
		session,
	)

	updateProduct.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	fieldsToUpdateProduct := &model.InternalUpdateProduct{
		ID: updateProduct.ID,
	}
	jsonExisting, _ := json.Marshal(existingProduct.ProposedChanges)
	json.Unmarshal(jsonExisting, &fieldsToUpdateProduct.ProposedChanges)

	var updateProductMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateProduct)
	json.Unmarshal(jsonUpdate, &updateProductMap)

	approveProdTrx.mapProcessorUtility.RemoveNil(updateProductMap)

	jsonUpdate, _ = json.Marshal(updateProductMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateProduct.ProposedChanges)

	if updateProduct.ProposalStatus != nil {
		if *updateProduct.ProposalStatus == model.EntityProposalStatusApproved {
			jsonTemp, _ := json.Marshal(fieldsToUpdateProduct.ProposedChanges)
			json.Unmarshal(jsonTemp, fieldsToUpdateProduct)
		}
	}

	updatedProduct, err := approveProdTrx.productDataSource.GetMongoDataSource().Update(
		fieldsToUpdateProduct.ID,
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
