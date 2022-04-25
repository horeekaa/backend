package supplyorderitemdomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasepurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type approveUpdateSupplyOrderItemTransactionComponent struct {
	supplyOrderItemDataSource              databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource
	purchaseOrderToSupplyDataSource        databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource
	approveUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ApproveUpdateDescriptivePhotoTransactionComponent
	loggingDataSource                      databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility                    coreutilityinterfaces.MapProcessorUtility
	pathIdentity                           string
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
		pathIdentity:                           "ApproveUpdateSupplyOrderItemComponent",
	}, nil
}

func (approveSupplyOrderItemTrx *approveUpdateSupplyOrderItemTransactionComponent) PreTransaction(
	input *model.InternalUpdateSupplyOrderItem,
) (*model.InternalUpdateSupplyOrderItem, error) {
	return input, nil
}

func (approveSupplyOrderItemTrx *approveUpdateSupplyOrderItemTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdateSupplyOrderItem,
) (*model.SupplyOrderItem, error) {
	updateSupplyOrderItem := &model.DatabaseUpdateSupplyOrderItem{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updateSupplyOrderItem)

	existingSupplyOrderItem, err := approveSupplyOrderItemTrx.supplyOrderItemDataSource.GetMongoDataSource().FindByID(
		updateSupplyOrderItem.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveSupplyOrderItemTrx.pathIdentity,
			err,
		)
	}
	if existingSupplyOrderItem.ProposedChanges.ProposalStatus != model.EntityProposalStatusProposed {
		return existingSupplyOrderItem, nil
	}

	previousLog, err := approveSupplyOrderItemTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingSupplyOrderItem.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveSupplyOrderItemTrx.pathIdentity,
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
	jsonTemp, _ = json.Marshal(
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
			approveSupplyOrderItemTrx.pathIdentity,
			err,
		)
	}

	updateSupplyOrderItem.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	currentTime := time.Now().UTC()
	updateSupplyOrderItem.UpdatedAt = &currentTime

	if updateSupplyOrderItem.SupplyOrderItemReturn != nil {
		if existingSupplyOrderItem.SupplyOrderItemReturn == nil {
			updateSupplyOrderItem.SupplyOrderItemReturn.CreatedAt = &currentTime
		}
		updateSupplyOrderItem.SupplyOrderItemReturn.UpdatedAt = &currentTime
	}

	fieldsToUpdateSupplyOrderItem := &model.DatabaseUpdateSupplyOrderItem{
		ID: updateSupplyOrderItem.ID,
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
				if *fieldsToUpdateSupplyOrderItem.ProposedChanges.QuantityAccepted > 0 {
					existingPOToSupply, err := approveSupplyOrderItemTrx.purchaseOrderToSupplyDataSource.GetMongoDataSource().FindByID(
						fieldsToUpdateSupplyOrderItem.ProposedChanges.PurchaseOrderToSupply.ID,
						session,
					)
					if err != nil {
						return nil, horeekaacoreexceptiontofailure.ConvertException(
							approveSupplyOrderItemTrx.pathIdentity,
							err,
						)
					}

					quantityFulfilled := existingPOToSupply.QuantityFulfilled +
						(existingSupplyOrderItem.ProposedChanges.QuantityAccepted -
							existingSupplyOrderItem.QuantityAccepted)
					poToSupplyToUpdate := &model.DatabaseUpdatePurchaseOrderToSupply{
						QuantityFulfilled: &quantityFulfilled,
					}
					poToSupplyToUpdate.Status = func(s model.PurchaseOrderToSupplyStatus) *model.PurchaseOrderToSupplyStatus { return &s }(model.PurchaseOrderToSupplyStatusOpen)
					if quantityFulfilled >= existingPOToSupply.QuantityRequested {
						poToSupplyToUpdate.Status = func(s model.PurchaseOrderToSupplyStatus) *model.PurchaseOrderToSupplyStatus { return &s }(model.PurchaseOrderToSupplyStatusFulfilled)
					}

					if !funk.Contains(
						funk.Map(existingPOToSupply.SupplyOrderItems,
							func(m *model.SupplyOrderItem) string { return m.ID.Hex() },
						),
						existingSupplyOrderItem.ID.Hex(),
					) {
						poToSupplyToUpdate.SupplyOrderItems = append(
							funk.Map(existingPOToSupply.SupplyOrderItems,
								func(m *model.SupplyOrderItem) *model.ObjectIDOnly {
									return &model.ObjectIDOnly{
										ID: &m.ID,
									}
								},
							).([]*model.ObjectIDOnly),
							&model.ObjectIDOnly{
								ID: &existingSupplyOrderItem.ID,
							},
						)
					}
					_, err = approveSupplyOrderItemTrx.purchaseOrderToSupplyDataSource.GetMongoDataSource().Update(
						map[string]interface{}{
							"_id": existingPOToSupply.ID,
						},
						poToSupplyToUpdate,
						session,
					)
					if err != nil {
						return nil, horeekaacoreexceptiontofailure.ConvertException(
							approveSupplyOrderItemTrx.pathIdentity,
							err,
						)
					}
				}
			}
		}

		if existingSupplyOrderItem.ProposedChanges.SupplyOrderItemReturn != nil {
			for _, updateDescriptivePhoto := range existingSupplyOrderItem.ProposedChanges.SupplyOrderItemReturn.Photos {
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
					return nil, err
				}
			}
		}
		for _, updateDescriptivePhoto := range existingSupplyOrderItem.ProposedChanges.Photos {
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
				return nil, err
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
			approveSupplyOrderItemTrx.pathIdentity,
			err,
		)
	}

	return updatedSupplyOrderItem, nil
}
