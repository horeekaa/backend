package supplyorderitemdomainrepositories

import (
	"encoding/json"
	"strings"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasepurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	supplyorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type proposeUpdateSupplyOrderItemTransactionComponent struct {
	supplyOrderItemDataSource              databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource
	loggingDataSource                      databaseloggingdatasourceinterfaces.LoggingDataSource
	purchaseOrderToSupplyDataSource        databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource
	createDescriptivePhotoComponent        descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent
	proposeUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ProposeUpdateDescriptivePhotoTransactionComponent
	supplyOrderItemLoader                  supplyorderitemdomainrepositoryutilityinterfaces.SupplyOrderItemLoader
	mapProcessorUtility                    coreutilityinterfaces.MapProcessorUtility
}

func NewProposeUpdateSupplyOrderItemTransactionComponent(
	supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
	proposeUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ProposeUpdateDescriptivePhotoTransactionComponent,
	supplyOrderItemLoader supplyorderitemdomainrepositoryutilityinterfaces.SupplyOrderItemLoader,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (supplyorderitemdomainrepositoryinterfaces.ProposeUpdateSupplyOrderItemTransactionComponent, error) {
	return &proposeUpdateSupplyOrderItemTransactionComponent{
		supplyOrderItemDataSource:              supplyOrderItemDataSource,
		loggingDataSource:                      loggingDataSource,
		purchaseOrderToSupplyDataSource:        purchaseOrderToSupplyDataSource,
		createDescriptivePhotoComponent:        createDescriptivePhotoComponent,
		proposeUpdateDescriptivePhotoComponent: proposeUpdateDescriptivePhotoComponent,
		supplyOrderItemLoader:                  supplyOrderItemLoader,
		mapProcessorUtility:                    mapProcessorUtility,
	}, nil
}

func (updateSupplyOrderItemTrx *proposeUpdateSupplyOrderItemTransactionComponent) PreTransaction(
	input *model.InternalUpdateSupplyOrderItem,
) (*model.InternalUpdateSupplyOrderItem, error) {
	return input, nil
}

func (updateSupplyOrderItemTrx *proposeUpdateSupplyOrderItemTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updateSupplyOrderItem *model.InternalUpdateSupplyOrderItem,
) (*model.SupplyOrderItem, error) {
	existingSupplyOrderItem, err := updateSupplyOrderItemTrx.supplyOrderItemDataSource.GetMongoDataSource().FindByID(
		*updateSupplyOrderItem.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/proposeUpdateSupplyOrderItemComponent",
			err,
		)
	}

	if updateSupplyOrderItem.Photos != nil {
		savedPhotos := existingSupplyOrderItem.Photos
		for _, photoToUpdate := range updateSupplyOrderItem.Photos {
			if photoToUpdate.ID != nil {
				if !funk.Contains(
					existingSupplyOrderItem.Photos,
					func(dp *model.DescriptivePhoto) bool {
						return dp.ID == *photoToUpdate.ID
					},
				) {
					continue
				}

				photoToUpdate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*updateSupplyOrderItem.ProposalStatus)
				photoToUpdate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*updateSupplyOrderItem.SubmittingAccount)
				_, err := updateSupplyOrderItemTrx.proposeUpdateDescriptivePhotoComponent.TransactionBody(
					session,
					photoToUpdate,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/updateSupplyOrderItem",
						err,
					)
				}

				if photoToUpdate.IsActive != nil {
					if !*photoToUpdate.IsActive {
						index := funk.IndexOf(
							savedPhotos,
							func(dp *model.DescriptivePhoto) bool {
								return dp.ID == *photoToUpdate.ID
							},
						)
						if index > -1 {
							savedPhotos = append(savedPhotos[:index], savedPhotos[index+1:]...)
						}
					}
				}
				continue
			}
			photoToCreate := &model.InternalCreateDescriptivePhoto{}
			jsonTemp, _ := json.Marshal(photoToUpdate)
			json.Unmarshal(jsonTemp, photoToCreate)
			photoToCreate.Category = model.DescriptivePhotoCategorySupplyOrderItemOnPickup
			photoToCreate.Object = &model.ObjectIDOnly{
				ID: &existingSupplyOrderItem.ID,
			}
			photoToCreate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*updateSupplyOrderItem.ProposalStatus)
			photoToCreate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*updateSupplyOrderItem.SubmittingAccount)
			if photoToUpdate.Photo != nil {
				photoToCreate.Photo.File = photoToUpdate.Photo.File
			}
			createdDescriptivePhoto, err := updateSupplyOrderItemTrx.createDescriptivePhotoComponent.TransactionBody(
				session,
				photoToCreate,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/updateSupplyOrderItem",
					err,
				)
			}

			savedPhotos = append(savedPhotos, createdDescriptivePhoto)
		}
		jsonTemp, _ := json.Marshal(
			map[string]interface{}{
				"Photos": savedPhotos,
			},
		)
		json.Unmarshal(jsonTemp, updateSupplyOrderItem)
	}

	if updateSupplyOrderItem.PickUpDetail != nil {
		savedPhotos := existingSupplyOrderItem.PickUpDetail.Photos
		for _, photoToUpdate := range updateSupplyOrderItem.PickUpDetail.Photos {
			if photoToUpdate.ID != nil {
				if !funk.Contains(
					existingSupplyOrderItem.PickUpDetail.Photos,
					func(dp *model.DescriptivePhoto) bool {
						return dp.ID == *photoToUpdate.ID
					},
				) {
					continue
				}

				photoToUpdate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*updateSupplyOrderItem.ProposalStatus)
				photoToUpdate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*updateSupplyOrderItem.SubmittingAccount)
				_, err := updateSupplyOrderItemTrx.proposeUpdateDescriptivePhotoComponent.TransactionBody(
					session,
					photoToUpdate,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/updateSupplyOrderItem",
						err,
					)
				}

				if photoToUpdate.IsActive != nil {
					if !*photoToUpdate.IsActive {
						index := funk.IndexOf(
							savedPhotos,
							func(dp *model.DescriptivePhoto) bool {
								return dp.ID == *photoToUpdate.ID
							},
						)
						if index > -1 {
							savedPhotos = append(savedPhotos[:index], savedPhotos[index+1:]...)
						}
					}
				}
				continue
			}
			photoToCreate := &model.InternalCreateDescriptivePhoto{}
			jsonTemp, _ := json.Marshal(photoToUpdate)
			json.Unmarshal(jsonTemp, photoToCreate)
			photoToCreate.Category = model.DescriptivePhotoCategorySupplyOrderItem
			photoToCreate.Object = &model.ObjectIDOnly{
				ID: &existingSupplyOrderItem.ID,
			}
			photoToCreate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*updateSupplyOrderItem.ProposalStatus)
			photoToCreate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*updateSupplyOrderItem.SubmittingAccount)
			if photoToUpdate.Photo != nil {
				photoToCreate.Photo.File = photoToUpdate.Photo.File
			}
			createdDescriptivePhoto, err := updateSupplyOrderItemTrx.createDescriptivePhotoComponent.TransactionBody(
				session,
				photoToCreate,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/updateSupplyOrderItem",
					err,
				)
			}

			savedPhotos = append(savedPhotos, createdDescriptivePhoto)
		}
		jsonTemp, _ := json.Marshal(
			map[string]interface{}{
				"PickUpDetail": map[string]interface{}{
					"Photos": savedPhotos,
				},
			},
		)
		json.Unmarshal(jsonTemp, updateSupplyOrderItem)

		if updateSupplyOrderItem.PickUpDetail.Courier != nil {
			generatedObjectID := updateSupplyOrderItemTrx.supplyOrderItemDataSource.GetMongoDataSource().GenerateObjectID()
			loc, _ := time.LoadLocation("Asia/Bangkok")
			splittedId := strings.Split(generatedObjectID.Hex(), "")
			updateSupplyOrderItem.PickUpDetail.PublicID = func(s ...string) *string { joinedString := strings.Join(s, "/"); return &joinedString }(
				"PK",
				time.Now().In(loc).Format("20060102"),
				strings.ToUpper(
					strings.Join(
						splittedId[len(splittedId)-4:],
						"",
					),
				),
			)
		}

		if updateSupplyOrderItem.PickUpDetail.CourierResponded != nil {
			if *updateSupplyOrderItem.PickUpDetail.CourierResponded {
				updateSupplyOrderItem.PickUpDetail.Status = func(m model.PickUpStatus) *model.PickUpStatus {
					return &m
				}(model.PickUpStatusDriverAssigned)
			}
		}

		if updateSupplyOrderItem.PickUpDetail.StartedPickingUp != nil {
			if *updateSupplyOrderItem.PickUpDetail.StartedPickingUp {
				currentTime := time.Now().UTC()
				updateSupplyOrderItem.PickUpDetail.StartPickUpTime = &currentTime
				updateSupplyOrderItem.PickUpDetail.Status = func(m model.PickUpStatus) *model.PickUpStatus {
					return &m
				}(model.PickUpStatusPickingUp)
			}
		}

		if updateSupplyOrderItem.PickUpDetail.FinishedPickingUp != nil {
			if *updateSupplyOrderItem.PickUpDetail.FinishedPickingUp {
				currentTime := time.Now().UTC()
				updateSupplyOrderItem.PickUpDetail.FinishPickUpTime = &currentTime
				updateSupplyOrderItem.PickUpDetail.Status = func(m model.PickUpStatus) *model.PickUpStatus {
					return &m
				}(model.PickUpStatusPickedUp)
			}
		}
	}

	_, err = updateSupplyOrderItemTrx.supplyOrderItemLoader.TransactionBody(
		session,
		updateSupplyOrderItem.PurchaseOrderToSupply,
		updateSupplyOrderItem.PickUpDetail,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/proposeUpdateSupplyOrderItemComponent",
			err,
		)
	}

	unitPrice := existingSupplyOrderItem.UnitPrice
	if updateSupplyOrderItem.UnitPrice != nil {
		unitPrice = *updateSupplyOrderItem.UnitPrice
	}

	quantity := existingSupplyOrderItem.QuantityOffered
	if updateSupplyOrderItem.QuantityOffered != nil {
		quantity = *updateSupplyOrderItem.QuantityOffered
	}

	quantityAccepted := existingSupplyOrderItem.QuantityAccepted
	if updateSupplyOrderItem.QuantityAccepted != nil {
		quantityAccepted = *updateSupplyOrderItem.QuantityAccepted
	}
	subTotal := quantity * unitPrice
	if quantityAccepted > 0 {
		subTotal = quantityAccepted * unitPrice
	}
	updateSupplyOrderItem.SubTotal = &subTotal

	subTotalReturn := 0
	if existingSupplyOrderItem.SupplyOrderItemReturn != nil {
		subTotalReturn = existingSupplyOrderItem.SupplyOrderItemReturn.SubTotal
	}
	if updateSupplyOrderItem.SupplyOrderItemReturn != nil {
		subTotalReturn = *updateSupplyOrderItem.SupplyOrderItemReturn.Quantity * unitPrice
		if subTotalReturn > subTotal {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.SOReturnAmountExceedFulfilledAmount,
				"/updateSupplyOrderItem",
				nil,
			)
		}
		updateSupplyOrderItem.SupplyOrderItemReturn.SubTotal = &subTotalReturn
	}
	salesAmount := subTotal - subTotalReturn
	updateSupplyOrderItem.SalesAmount = &salesAmount

	if updateSupplyOrderItem.QuantityAccepted != nil {
		if *updateSupplyOrderItem.QuantityAccepted > 0 {
			existingPOToSupply, err := updateSupplyOrderItemTrx.purchaseOrderToSupplyDataSource.GetMongoDataSource().FindByID(
				existingSupplyOrderItem.PurchaseOrderToSupply.ID,
				session,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/proposeUpdateSupplyOrderItem",
					err,
				)
			}
			if *updateSupplyOrderItem.QuantityAccepted < existingSupplyOrderItem.QuantityOffered {
				updateSupplyOrderItem.Status = func(m model.SupplyOrderItemStatus) *model.SupplyOrderItemStatus {
					return &m
				}(model.SupplyOrderItemStatusPartiallyAccepted)
			} else {
				updateSupplyOrderItem.Status = func(m model.SupplyOrderItemStatus) *model.SupplyOrderItemStatus {
					return &m
				}(model.SupplyOrderItemStatusAccepted)
				updateSupplyOrderItem.PartnerAgreed = func(b bool) *bool {
					return &b
				}(true)

				quantityFulfilled := existingPOToSupply.QuantityFulfilled +
					(*updateSupplyOrderItem.QuantityAccepted - existingSupplyOrderItem.QuantityAccepted)

				poToSupplyToUpdate := &model.DatabaseUpdatePurchaseOrderToSupply{
					QuantityFulfilled: &quantityFulfilled,
				}
				poToSupplyToUpdate.Status = func(s model.PurchaseOrderToSupplyStatus) *model.PurchaseOrderToSupplyStatus {
					return &s
				}(model.PurchaseOrderToSupplyStatusOpen)
				if quantityFulfilled >= existingPOToSupply.QuantityRequested {
					poToSupplyToUpdate.Status = func(s model.PurchaseOrderToSupplyStatus) *model.PurchaseOrderToSupplyStatus {
						return &s
					}(model.PurchaseOrderToSupplyStatusFulfilled)
				}
				_, err = updateSupplyOrderItemTrx.purchaseOrderToSupplyDataSource.GetMongoDataSource().Update(
					map[string]interface{}{
						"_id": existingPOToSupply.ID,
					},
					poToSupplyToUpdate,
					session,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/proposeUpdateSupplyOrderItem",
						err,
					)
				}
			}
		}
	}

	if updateSupplyOrderItem.SupplyOrderItemReturn != nil {
		if existingSupplyOrderItem.SupplyOrderItemReturn == nil {
			generatedObjectID := updateSupplyOrderItemTrx.supplyOrderItemDataSource.GetMongoDataSource().GenerateObjectID()
			loc, _ := time.LoadLocation("Asia/Bangkok")
			splittedId := strings.Split(generatedObjectID.Hex(), "")
			updateSupplyOrderItem.SupplyOrderItemReturn.PublicID = func(s ...string) *string { joinedString := strings.Join(s, "/"); return &joinedString }(
				"ISR",
				time.Now().In(loc).Format("20060102"),
				strings.ToUpper(
					strings.Join(
						splittedId[len(splittedId)-4:],
						"",
					),
				),
			)
		}
	}

	newDocumentJson, _ := json.Marshal(*updateSupplyOrderItem)
	oldDocumentJson, _ := json.Marshal(*existingSupplyOrderItem)
	loggingOutput, err := updateSupplyOrderItemTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "SupplyOrderItem",
			Document: &model.ObjectIDOnly{
				ID: &existingSupplyOrderItem.ID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			OldDocumentJSON: func(s string) *string { return &s }(string(oldDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updateSupplyOrderItem.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updateSupplyOrderItem.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/proposeUpdateSupplyOrderItemComponent",
			err,
		)
	}
	updateSupplyOrderItem.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	fieldsToUpdatesupplyOrderItem := &model.DatabaseUpdateSupplyOrderItem{
		ID: *updateSupplyOrderItem.ID,
	}
	jsonExisting, _ := json.Marshal(existingSupplyOrderItem)
	json.Unmarshal(jsonExisting, &fieldsToUpdatesupplyOrderItem.ProposedChanges)

	var updatesupplyOrderItemMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updateSupplyOrderItem)
	json.Unmarshal(jsonUpdate, &updatesupplyOrderItemMap)

	updateSupplyOrderItemTrx.mapProcessorUtility.RemoveNil(updatesupplyOrderItemMap)

	jsonUpdate, _ = json.Marshal(updatesupplyOrderItemMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdatesupplyOrderItem.ProposedChanges)

	if updateSupplyOrderItem.ProposalStatus != nil {
		fieldsToUpdatesupplyOrderItem.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: updateSupplyOrderItem.SubmittingAccount.ID,
		}
		if *updateSupplyOrderItem.ProposalStatus == model.EntityProposalStatusApproved {
			json.Unmarshal(jsonUpdate, fieldsToUpdatesupplyOrderItem)
		}
	}

	updatedSupplyOrderItem, err := updateSupplyOrderItemTrx.supplyOrderItemDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdatesupplyOrderItem.ID,
		},
		fieldsToUpdatesupplyOrderItem,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/proposeUpdateSupplyOrderItemComponent",
			err,
		)
	}

	return updatedSupplyOrderItem, nil
}
