package productvariantdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databaseproductvariantdatasourceinterfaces "github.com/horeekaa/backend/features/productVariants/data/dataSources/databases/interfaces/sources"
	productvariantdomainrepositoryinterfaces "github.com/horeekaa/backend/features/productVariants/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateProductVariantTransactionComponent struct {
	productVariantDataSource        databaseproductvariantdatasourceinterfaces.ProductVariantDataSource
	loggingDataSource               databaseloggingdatasourceinterfaces.LoggingDataSource
	updateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.UpdateDescriptivePhotoTransactionComponent
	mapProcessorUtility             coreutilityinterfaces.MapProcessorUtility
}

func NewApproveUpdateProductVariantTransactionComponent(
	productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	updateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.UpdateDescriptivePhotoTransactionComponent,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (productvariantdomainrepositoryinterfaces.ApproveUpdateProductVariantTransactionComponent, error) {
	return &approveUpdateProductVariantTransactionComponent{
		productVariantDataSource:        productVariantDataSource,
		loggingDataSource:               loggingDataSource,
		updateDescriptivePhotoComponent: updateDescriptivePhotoComponent,
		mapProcessorUtility:             mapProcessorUtility,
	}, nil
}

func (approveProdVarTrx *approveUpdateProductVariantTransactionComponent) PreTransaction(
	input *model.InternalUpdateProductVariant,
) (*model.InternalUpdateProductVariant, error) {
	return input, nil
}

func (approveProdVarTrx *approveUpdateProductVariantTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateProductVariant *model.InternalUpdateProductVariant,
) (*model.ProductVariant, error) {
	existingProductVariant, err := approveProdVarTrx.productVariantDataSource.GetMongoDataSource().FindByID(
		*updateProductVariant.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateProductVariant",
			err,
		)
	}
	if existingProductVariant.ProposedChanges.ProposalStatus == model.EntityProposalStatusApproved {
		return existingProductVariant, nil
	}

	previousLog, err := approveProdVarTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingProductVariant.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateProductVariant",
			err,
		)
	}

	logToCreate := &model.CreateLogging{
		Collection: previousLog.Collection,
		CreatedByAccount: &model.ObjectIDOnly{
			ID: updateProductVariant.RecentApprovingAccount.ID,
		},
		Activity:       previousLog.Activity,
		ProposalStatus: *updateProductVariant.ProposalStatus,
	}
	jsonTemp, _ := json.Marshal(
		map[string]interface{}{
			"NewDocumentJSON": previousLog.NewDocumentJSON,
			"OldDocumentJSON": previousLog.OldDocumentJSON,
		},
	)
	json.Unmarshal(jsonTemp, logToCreate)

	createdLog, err := approveProdVarTrx.loggingDataSource.GetMongoDataSource().Create(
		logToCreate,
		session,
	)

	updateProductVariant.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	fieldsToUpdateProductVariant := &model.DatabaseUpdateProductVariant{
		ID: *updateProductVariant.ID,
	}
	jsonExisting, _ := json.Marshal(existingProductVariant.ProposedChanges)
	json.Unmarshal(jsonExisting, &fieldsToUpdateProductVariant.ProposedChanges)

	var updateProductVariantMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateProductVariant)
	json.Unmarshal(jsonUpdate, &updateProductVariantMap)

	approveProdVarTrx.mapProcessorUtility.RemoveNil(updateProductVariantMap)

	jsonUpdate, _ = json.Marshal(updateProductVariantMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateProductVariant.ProposedChanges)

	if updateProductVariant.ProposalStatus != nil {
		if *updateProductVariant.ProposalStatus == model.EntityProposalStatusApproved {
			jsonUpdate, _ := json.Marshal(fieldsToUpdateProductVariant.ProposedChanges)
			json.Unmarshal(jsonUpdate, fieldsToUpdateProductVariant)
		}

		if *updateProductVariant.ProposalStatus == model.EntityProposalStatusRejected &&
			existingProductVariant.ProposalStatus == model.EntityProposalStatusProposed {
			if existingProductVariant.Photo != nil {
				_, err = approveProdVarTrx.updateDescriptivePhotoComponent.TransactionBody(
					session,
					&model.InternalUpdateDescriptivePhoto{
						ID:       &existingProductVariant.Photo.ID,
						IsActive: func(b bool) *bool { return &b }(false),
					},
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/updateProductVariant",
						err,
					)
				}
			}
		}
	}

	updatedProductVariant, err := approveProdVarTrx.productVariantDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateProductVariant.ID,
		},
		fieldsToUpdateProductVariant,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateProductVariant",
			err,
		)
	}

	return updatedProductVariant, nil
}
