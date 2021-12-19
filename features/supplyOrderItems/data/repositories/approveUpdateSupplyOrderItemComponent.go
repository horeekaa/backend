package supplyorderitemdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasepurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateSupplyOrderItemTransactionComponent struct {
	supplyOrderItemDataSource              databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource
	purchaseOrderToSupplyDataSource        databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource
	approveUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent
	loggingDataSource                      databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility                    coreutilityinterfaces.MapProcessorUtility
}

func NewApproveUpdateSupplyOrderItemTransactionComponent(
	supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource,
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource,
	approveUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (supplyorderitemdomainrepositoryinterfaces.ApproveUpdateSupplyOrderItemTransactionComponent, error) {
	return &approveUpdateSupplyOrderItemTransactionComponent{
		supplyOrderItemDataSource:              supplyOrderItemDataSource,
		purchaseOrderToSupplyDataSource:        purchaseOrderToSupplyDataSource,
		approveUpdateDescriptivePhotoComponent: approveUpdateDescriptivePhotoComponent,
		loggingDataSource:                      loggingDataSource,
		mapProcessorUtility:                    mapProcessorUtility,
	}, nil
}

func (approveSupplyOrderItemTrx *approveUpdateSupplyOrderItemTransactionComponent) PreTransaction(
	input *model.InternalUpdateSupplyOrderItem,
) (*model.InternalUpdateSupplyOrderItem, error) {
	return input, nil
}

func (approveSupplyOrderItemTrx *approveUpdateSupplyOrderItemTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateSupplyOrderItem *model.InternalUpdateSupplyOrderItem,
) (*model.SupplyOrderItem, error) {
	existingSupplyOrderItem, err := approveSupplyOrderItemTrx.supplyOrderItemDataSource.GetMongoDataSource().FindByID(
		*updateSupplyOrderItem.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/approveUpdateSupplyOrderItem",
			err,
		)
	}
	if existingSupplyOrderItem.ProposedChanges.ProposalStatus == model.EntityProposalStatusApproved {
		return existingSupplyOrderItem, nil
	}

	previousLog, err := approveSupplyOrderItemTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingSupplyOrderItem.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/approveUpdateSupplyOrderItem",
			err,
		)
	}

	logToCreate := &model.CreateLogging{
		Collection: previousLog.Collection,
		CreatedByAccount: &model.ObjectIDOnly{
			ID: updateSupplyOrderItem.RecentApprovingAccount.ID,
		},
		Activity:       previousLog.Activity,
		ProposalStatus: *updateSupplyOrderItem.ProposalStatus,
	}
	jsonTemp, _ := json.Marshal(
		map[string]interface{}{
			"NewDocumentJSON": previousLog.NewDocumentJSON,
			"OldDocumentJSON": previousLog.OldDocumentJSON,
		},
	)
	json.Unmarshal(jsonTemp, logToCreate)

	createdLog, err := approveSupplyOrderItemTrx.loggingDataSource.GetMongoDataSource().Create(
		logToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/approveUpdateSupplyOrderItem",
			err,
		)
	}

	updateSupplyOrderItem.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	fieldsToUpdateSupplyOrderItem := &model.DatabaseUpdateSupplyOrderItem{
		ID: *updateSupplyOrderItem.ID,
	}
	jsonExisting, _ := json.Marshal(existingSupplyOrderItem.ProposedChanges)
	json.Unmarshal(jsonExisting, &fieldsToUpdateSupplyOrderItem.ProposedChanges)

	var updatesupplyOrderItemMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateSupplyOrderItem)
	json.Unmarshal(jsonUpdate, &updatesupplyOrderItemMap)

	approveSupplyOrderItemTrx.mapProcessorUtility.RemoveNil(updatesupplyOrderItemMap)

	jsonUpdate, _ = json.Marshal(updatesupplyOrderItemMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdateSupplyOrderItem.ProposedChanges)

	if updateSupplyOrderItem.ProposalStatus != nil {
		if *updateSupplyOrderItem.ProposalStatus == model.EntityProposalStatusApproved {
			jsonUpdate, _ := json.Marshal(fieldsToUpdateSupplyOrderItem.ProposedChanges)
			json.Unmarshal(jsonUpdate, fieldsToUpdateSupplyOrderItem)

			if *fieldsToUpdateSupplyOrderItem.ProposedChanges.PartnerAgreed {
				existingPOToSupply, err := approveSupplyOrderItemTrx.purchaseOrderToSupplyDataSource.GetMongoDataSource().FindByID(
					fieldsToUpdateSupplyOrderItem.ProposedChanges.PurchaseOrderToSupply.ID,
					session,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/approveUpdateSupplyOrderItem",
						err,
					)
				}

				quantityFulfilled := existingPOToSupply.QuantityFulfilled +
					(*fieldsToUpdateSupplyOrderItem.ProposedChanges.QuantityAccepted -
						existingSupplyOrderItem.QuantityAccepted)
				poToSupplyToUpdate := &model.DatabaseUpdatePurchaseOrderToSupply{
					QuantityFulfilled: &quantityFulfilled,
				}
				poToSupplyToUpdate.Status = func(s model.PurchaseOrderToSupplyStatus) *model.PurchaseOrderToSupplyStatus { return &s }(model.PurchaseOrderToSupplyStatusOpen)
				if quantityFulfilled >= existingPOToSupply.QuantityRequested {
					poToSupplyToUpdate.Status = func(s model.PurchaseOrderToSupplyStatus) *model.PurchaseOrderToSupplyStatus { return &s }(model.PurchaseOrderToSupplyStatusFulfilled)
				}
				_, err = approveSupplyOrderItemTrx.purchaseOrderToSupplyDataSource.GetMongoDataSource().Update(
					map[string]interface{}{
						"_id": fieldsToUpdateSupplyOrderItem.ProposedChanges.PurchaseOrderToSupply.ID,
					},
					poToSupplyToUpdate,
					session,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/approveUpdateSupplyOrderItem",
						err,
					)
				}
			}
		}

		if existingSupplyOrderItem.PickUpDetail != nil {
			for _, updateDescriptivePhoto := range existingSupplyOrderItem.PickUpDetail.Photos {
				updateDescriptivePhoto := &model.InternalUpdateDescriptivePhoto{
					ID: &updateDescriptivePhoto.ID,
				}
				updateDescriptivePhoto.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*updateSupplyOrderItem.RecentApprovingAccount)
				updateDescriptivePhoto.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*updateSupplyOrderItem.ProposalStatus)

				_, err := approveSupplyOrderItemTrx.approveUpdateDescriptivePhotoComponent.TransactionBody(
					session,
					updateDescriptivePhoto,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/approveUpdateSupplyOrderItem",
						err,
					)
				}
			}
		}
	}

	updatedSupplyOrderItem, err := approveSupplyOrderItemTrx.supplyOrderItemDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdateSupplyOrderItem.ID,
		},
		fieldsToUpdateSupplyOrderItem,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/approveUpdateSupplyOrderItem",
			err,
		)
	}

	return updatedSupplyOrderItem, nil
}
