package purchaseorderitemdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdatePurchaseOrderItemTransactionComponent struct {
	purchaseOrderItemDataSource            databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource
	loggingDataSource                      databaseloggingdatasourceinterfaces.LoggingDataSource
	approveUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent
	mapProcessorUtility                    coreutilityinterfaces.MapProcessorUtility
}

func NewApproveUpdatePurchaseOrderItemTransactionComponent(
	purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	approveUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (purchaseorderitemdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderItemTransactionComponent, error) {
	return &approveUpdatePurchaseOrderItemTransactionComponent{
		purchaseOrderItemDataSource:            purchaseOrderItemDataSource,
		loggingDataSource:                      loggingDataSource,
		approveUpdateDescriptivePhotoComponent: approveUpdateDescriptivePhotoComponent,
		mapProcessorUtility:                    mapProcessorUtility,
	}, nil
}

func (approvePOItemTrx *approveUpdatePurchaseOrderItemTransactionComponent) PreTransaction(
	input *model.InternalUpdatePurchaseOrderItem,
) (*model.InternalUpdatePurchaseOrderItem, error) {
	return input, nil
}

func (approvePOItemTrx *approveUpdatePurchaseOrderItemTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updatePurchaseOrderItem *model.InternalUpdatePurchaseOrderItem,
) (*model.PurchaseOrderItem, error) {
	existingPurchaseOrderItem, err := approvePOItemTrx.purchaseOrderItemDataSource.GetMongoDataSource().FindByID(
		*updatePurchaseOrderItem.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/approveUpdatePurchaseOrderItem",
			err,
		)
	}
	if existingPurchaseOrderItem.ProposedChanges.ProposalStatus == model.EntityProposalStatusApproved {
		return existingPurchaseOrderItem, nil
	}

	previousLog, err := approvePOItemTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingPurchaseOrderItem.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/approveUpdatePurchaseOrderItem",
			err,
		)
	}

	logToCreate := &model.CreateLogging{
		Collection: previousLog.Collection,
		CreatedByAccount: &model.ObjectIDOnly{
			ID: updatePurchaseOrderItem.RecentApprovingAccount.ID,
		},
		Activity:       previousLog.Activity,
		ProposalStatus: *updatePurchaseOrderItem.ProposalStatus,
	}
	jsonTemp, _ := json.Marshal(
		map[string]interface{}{
			"NewDocumentJSON": previousLog.NewDocumentJSON,
			"OldDocumentJSON": previousLog.OldDocumentJSON,
		},
	)
	json.Unmarshal(jsonTemp, logToCreate)

	createdLog, err := approvePOItemTrx.loggingDataSource.GetMongoDataSource().Create(
		logToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/approveUpdatePurchaseOrderItem",
			err,
		)
	}

	updatePurchaseOrderItem.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	fieldsToUpdatePurchaseOrderItem := &model.DatabaseUpdatePurchaseOrderItem{
		ID: *updatePurchaseOrderItem.ID,
	}
	jsonExisting, _ := json.Marshal(existingPurchaseOrderItem.ProposedChanges)
	json.Unmarshal(jsonExisting, &fieldsToUpdatePurchaseOrderItem.ProposedChanges)

	var updatePurchaseOrderItemMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updatePurchaseOrderItem)
	json.Unmarshal(jsonUpdate, &updatePurchaseOrderItemMap)

	approvePOItemTrx.mapProcessorUtility.RemoveNil(updatePurchaseOrderItemMap)

	jsonUpdate, _ = json.Marshal(updatePurchaseOrderItemMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdatePurchaseOrderItem.ProposedChanges)

	if updatePurchaseOrderItem.ProposalStatus != nil {
		if *updatePurchaseOrderItem.ProposalStatus == model.EntityProposalStatusApproved {
			jsonUpdate, _ := json.Marshal(fieldsToUpdatePurchaseOrderItem.ProposedChanges)
			json.Unmarshal(jsonUpdate, fieldsToUpdatePurchaseOrderItem)
		}

		if existingPurchaseOrderItem.DeliveryDetail != nil {
			for _, updateDescriptivePhoto := range existingPurchaseOrderItem.DeliveryDetail.Photos {
				updateDescriptivePhoto := &model.InternalUpdateDescriptivePhoto{
					ID: &updateDescriptivePhoto.ID,
				}
				updateDescriptivePhoto.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*updatePurchaseOrderItem.RecentApprovingAccount)
				updateDescriptivePhoto.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*updatePurchaseOrderItem.ProposalStatus)

				_, err := approvePOItemTrx.approveUpdateDescriptivePhotoComponent.TransactionBody(
					session,
					updateDescriptivePhoto,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/approveUpdatePurchaseOrderItem",
						err,
					)
				}
			}
		}
	}

	updatedPurchaseOrderItem, err := approvePOItemTrx.purchaseOrderItemDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdatePurchaseOrderItem.ID,
		},
		fieldsToUpdatePurchaseOrderItem,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/approveUpdatePurchaseOrderItem",
			err,
		)
	}

	return updatedPurchaseOrderItem, nil
}
