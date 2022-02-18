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
	productVariantDataSource               databaseproductvariantdatasourceinterfaces.ProductVariantDataSource
	loggingDataSource                      databaseloggingdatasourceinterfaces.LoggingDataSource
	approveUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent
	mapProcessorUtility                    coreutilityinterfaces.MapProcessorUtility
}

func NewApproveUpdateProductVariantTransactionComponent(
	productVariantDataSource databaseproductvariantdatasourceinterfaces.ProductVariantDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	approveUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (productvariantdomainrepositoryinterfaces.ApproveUpdateProductVariantTransactionComponent, error) {
	return &approveUpdateProductVariantTransactionComponent{
		productVariantDataSource:               productVariantDataSource,
		loggingDataSource:                      loggingDataSource,
		approveUpdateDescriptivePhotoComponent: approveUpdateDescriptivePhotoComponent,
		mapProcessorUtility:                    mapProcessorUtility,
	}, nil
}

func (approveProdVarTrx *approveUpdateProductVariantTransactionComponent) PreTransaction(
	input *model.InternalUpdateProductVariant,
) (*model.InternalUpdateProductVariant, error) {
	return input, nil
}

func (approveProdVarTrx *approveUpdateProductVariantTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateProductVariant,
) (*model.ProductVariant, error) {
	updateProductVariant := &model.DatabaseUpdateProductVariant{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateProductVariant)

	existingProductVariant, err := approveProdVarTrx.productVariantDataSource.GetMongoDataSource().FindByID(
		updateProductVariant.ID,
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
	jsonTemp, _ = json.Marshal(
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
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updateProductVariant",
			err,
		)
	}

	updateProductVariant.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	fieldsToUpdateProductVariant := &model.DatabaseUpdateProductVariant{
		ID: updateProductVariant.ID,
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

		if existingProductVariant.ProposedChanges.Photo != nil {
			updateDescriptivePhoto := &model.InternalUpdateDescriptivePhoto{
				ID: &existingProductVariant.ProposedChanges.Photo.ID,
			}
			updateDescriptivePhoto.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*updateProductVariant.RecentApprovingAccount)
			updateDescriptivePhoto.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*updateProductVariant.ProposalStatus)

			_, err := approveProdVarTrx.approveUpdateDescriptivePhotoComponent.TransactionBody(
				session,
				updateDescriptivePhoto,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/updateProductVariant",
					err,
				)
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
