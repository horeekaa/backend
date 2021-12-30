package purchaseorderitemdomainrepositories

import (
	"encoding/json"
	"strings"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	descriptivephotodomainrepositoryinterfaces "github.com/horeekaa/backend/features/descriptivePhotos/domain/repositories"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasepurchaseOrderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	purchaseorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type proposeUpdatePurchaseOrderItemTransactionComponent struct {
	purchaseOrderItemDataSource            databasepurchaseOrderitemdatasourceinterfaces.PurchaseOrderItemDataSource
	loggingDataSource                      databaseloggingdatasourceinterfaces.LoggingDataSource
	createDescriptivePhotoComponent        descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent
	proposeUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ProposeUpdateDescriptivePhotoTransactionComponent
	purchaseOrderItemLoader                purchaseorderitemdomainrepositoryutilityinterfaces.PurchaseOrderItemLoader
	mapProcessorUtility                    coreutilityinterfaces.MapProcessorUtility
}

func NewProposeUpdatePurchaseOrderItemTransactionComponent(
	purchaseOrderItemDataSource databasepurchaseOrderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	createDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.CreateDescriptivePhotoTransactionComponent,
	proposeUpdateDescriptivePhotoComponent descriptivephotodomainrepositoryinterfaces.ProposeUpdateDescriptivePhotoTransactionComponent,
	purchaseOrderItemLoader purchaseorderitemdomainrepositoryutilityinterfaces.PurchaseOrderItemLoader,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (purchaseorderitemdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderItemTransactionComponent, error) {
	return &proposeUpdatePurchaseOrderItemTransactionComponent{
		purchaseOrderItemDataSource:            purchaseOrderItemDataSource,
		loggingDataSource:                      loggingDataSource,
		createDescriptivePhotoComponent:        createDescriptivePhotoComponent,
		proposeUpdateDescriptivePhotoComponent: proposeUpdateDescriptivePhotoComponent,
		purchaseOrderItemLoader:                purchaseOrderItemLoader,
		mapProcessorUtility:                    mapProcessorUtility,
	}, nil
}

func (updatePurchaseOrderItemTrx *proposeUpdatePurchaseOrderItemTransactionComponent) PreTransaction(
	input *model.InternalUpdatePurchaseOrderItem,
) (*model.InternalUpdatePurchaseOrderItem, error) {
	return input, nil
}

func (updatePurchaseOrderItemTrx *proposeUpdatePurchaseOrderItemTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updatePurchaseOrderItem *model.InternalUpdatePurchaseOrderItem,
) (*model.PurchaseOrderItem, error) {
	existingPurchaseOrderItem, err := updatePurchaseOrderItemTrx.purchaseOrderItemDataSource.GetMongoDataSource().FindByID(
		*updatePurchaseOrderItem.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/proposeUpdatePurchaseOrderItemComponent",
			err,
		)
	}

	if updatePurchaseOrderItem.DeliveryDetail != nil {
		savedPhotosAfterReceived := existingPurchaseOrderItem.DeliveryDetail.PhotosAfterReceived
		for _, photoToUpdate := range updatePurchaseOrderItem.DeliveryDetail.PhotosAfterReceived {
			if photoToUpdate.ID != nil {
				if !funk.Contains(
					existingPurchaseOrderItem.DeliveryDetail.PhotosAfterReceived,
					func(dp *model.DescriptivePhoto) bool {
						return dp.ID == *photoToUpdate.ID
					},
				) {
					continue
				}

				photoToUpdate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*updatePurchaseOrderItem.ProposalStatus)
				photoToUpdate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*updatePurchaseOrderItem.SubmittingAccount)
				_, err := updatePurchaseOrderItemTrx.proposeUpdateDescriptivePhotoComponent.TransactionBody(
					session,
					photoToUpdate,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/updatePurchaseOrderItem",
						err,
					)
				}

				if photoToUpdate.IsActive != nil {
					if !*photoToUpdate.IsActive {
						index := funk.IndexOf(
							savedPhotosAfterReceived,
							func(dp *model.DescriptivePhoto) bool {
								return dp.ID == *photoToUpdate.ID
							},
						)
						if index > -1 {
							savedPhotosAfterReceived = append(savedPhotosAfterReceived[:index], savedPhotosAfterReceived[index+1:]...)
						}
					}
				}
				continue
			}
			photoToCreate := &model.InternalCreateDescriptivePhoto{}
			jsonTemp, _ := json.Marshal(photoToUpdate)
			json.Unmarshal(jsonTemp, photoToCreate)
			photoToCreate.Category = model.DescriptivePhotoCategoryPurchaseOrderItemAfterReceived
			photoToCreate.Object = &model.ObjectIDOnly{
				ID: &existingPurchaseOrderItem.ID,
			}
			photoToCreate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*updatePurchaseOrderItem.ProposalStatus)
			photoToCreate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*updatePurchaseOrderItem.SubmittingAccount)
			if photoToUpdate.Photo != nil {
				photoToCreate.Photo.File = photoToUpdate.Photo.File
			}
			createdAfterReceivedPhoto, err := updatePurchaseOrderItemTrx.createDescriptivePhotoComponent.TransactionBody(
				session,
				photoToCreate,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/updatePurchaseOrderItem",
					err,
				)
			}

			savedPhotosAfterReceived = append(savedPhotosAfterReceived, createdAfterReceivedPhoto)
		}
		jsonTemp, _ := json.Marshal(
			map[string]interface{}{
				"DeliveryDetail": map[string]interface{}{
					"PhotosAfterReceived": savedPhotosAfterReceived,
				},
			},
		)
		json.Unmarshal(jsonTemp, updatePurchaseOrderItem)

		savedPhotos := existingPurchaseOrderItem.DeliveryDetail.Photos
		for _, photoToUpdate := range updatePurchaseOrderItem.DeliveryDetail.Photos {
			if photoToUpdate.ID != nil {
				if !funk.Contains(
					existingPurchaseOrderItem.DeliveryDetail.Photos,
					func(dp *model.DescriptivePhoto) bool {
						return dp.ID == *photoToUpdate.ID
					},
				) {
					continue
				}

				photoToUpdate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*updatePurchaseOrderItem.ProposalStatus)
				photoToUpdate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*updatePurchaseOrderItem.SubmittingAccount)
				_, err := updatePurchaseOrderItemTrx.proposeUpdateDescriptivePhotoComponent.TransactionBody(
					session,
					photoToUpdate,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/updatePurchaseOrderItem",
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
			photoToCreate.Category = model.DescriptivePhotoCategoryPurchaseOrderItem
			photoToCreate.Object = &model.ObjectIDOnly{
				ID: &existingPurchaseOrderItem.ID,
			}
			photoToCreate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*updatePurchaseOrderItem.ProposalStatus)
			photoToCreate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*updatePurchaseOrderItem.SubmittingAccount)
			if photoToUpdate.Photo != nil {
				photoToCreate.Photo.File = photoToUpdate.Photo.File
			}
			createdDescriptivePhoto, err := updatePurchaseOrderItemTrx.createDescriptivePhotoComponent.TransactionBody(
				session,
				photoToCreate,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/updatePurchaseOrderItem",
					err,
				)
			}

			savedPhotos = append(savedPhotos, createdDescriptivePhoto)
		}
		jsonTemp, _ = json.Marshal(
			map[string]interface{}{
				"DeliveryDetail": map[string]interface{}{
					"Photos": savedPhotos,
				},
			},
		)
		json.Unmarshal(jsonTemp, updatePurchaseOrderItem)

		if updatePurchaseOrderItem.DeliveryDetail.Courier != nil {
			generatedObjectID := updatePurchaseOrderItemTrx.purchaseOrderItemDataSource.GetMongoDataSource().GenerateObjectID()
			loc, _ := time.LoadLocation("Asia/Bangkok")
			splittedId := strings.Split(generatedObjectID.Hex(), "")
			updatePurchaseOrderItem.DeliveryDetail.PublicID = func(s ...string) *string { joinedString := strings.Join(s, "/"); return &joinedString }(
				"DV",
				time.Now().In(loc).Format("20060102"),
				strings.ToUpper(
					strings.Join(
						splittedId[len(splittedId)-4:],
						"",
					),
				),
			)
		}

		if updatePurchaseOrderItem.DeliveryDetail.CourierResponded != nil {
			if *updatePurchaseOrderItem.DeliveryDetail.CourierResponded {
				updatePurchaseOrderItem.DeliveryDetail.Status = func(m model.DeliveryStatus) *model.DeliveryStatus {
					return &m
				}(model.DeliveryStatusDriverAssigned)
			}
		}

		if updatePurchaseOrderItem.DeliveryDetail.StartedDelivering != nil {
			if *updatePurchaseOrderItem.DeliveryDetail.StartedDelivering {
				currentTime := time.Now().UTC()
				updatePurchaseOrderItem.DeliveryDetail.StartDeliveryTime = &currentTime
				updatePurchaseOrderItem.DeliveryDetail.Status = func(m model.DeliveryStatus) *model.DeliveryStatus {
					return &m
				}(model.DeliveryStatusDelivering)
			}
		}

		if updatePurchaseOrderItem.DeliveryDetail.FinishedDelivering != nil {
			if *updatePurchaseOrderItem.DeliveryDetail.FinishedDelivering {
				currentTime := time.Now().UTC()
				updatePurchaseOrderItem.DeliveryDetail.FinishDeliveryTime = &currentTime
				updatePurchaseOrderItem.DeliveryDetail.Status = func(m model.DeliveryStatus) *model.DeliveryStatus {
					return &m
				}(model.DeliveryStatusDelivered)
			}
		}
	}

	_, err = updatePurchaseOrderItemTrx.purchaseOrderItemLoader.TransactionBody(
		session,
		updatePurchaseOrderItem.MouItem,
		updatePurchaseOrderItem.ProductVariant,
		updatePurchaseOrderItem.DeliveryDetail,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/proposeUpdatePurchaseOrderItemComponent",
			err,
		)
	}

	if updatePurchaseOrderItem.ProductVariant != nil {
		updatePurchaseOrderItem.UnitPrice = &updatePurchaseOrderItem.ProductVariant.RetailPrice
		if existingPurchaseOrderItem.MouItem != nil {
			index := funk.IndexOf(
				existingPurchaseOrderItem.MouItem.AgreedProduct.Variants,
				func(pv *model.InternalAgreedProductVariantInput) bool {
					return pv.ID == updatePurchaseOrderItem.ProductVariant.ID
				},
			)
			if index > -1 {
				updatePurchaseOrderItem.UnitPrice = &existingPurchaseOrderItem.MouItem.AgreedProduct.Variants[index].RetailPrice
			}
		}
	}

	unitPrice := existingPurchaseOrderItem.UnitPrice
	if updatePurchaseOrderItem.UnitPrice != nil {
		unitPrice = *updatePurchaseOrderItem.UnitPrice
	}

	quantity := existingPurchaseOrderItem.Quantity
	if updatePurchaseOrderItem.Quantity != nil {
		quantity = *updatePurchaseOrderItem.Quantity
	}
	subTotal := quantity * unitPrice

	updatePurchaseOrderItem.SubTotal = &subTotal

	newDocumentJson, _ := json.Marshal(*updatePurchaseOrderItem)
	oldDocumentJson, _ := json.Marshal(*existingPurchaseOrderItem)
	loggingOutput, err := updatePurchaseOrderItemTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "PurchaseOrderItem",
			Document: &model.ObjectIDOnly{
				ID: &existingPurchaseOrderItem.ID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			OldDocumentJSON: func(s string) *string { return &s }(string(oldDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updatePurchaseOrderItem.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updatePurchaseOrderItem.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/proposeUpdatePurchaseOrderItemComponent",
			err,
		)
	}
	updatePurchaseOrderItem.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	fieldsToUpdatePurchaseOrderItem := &model.DatabaseUpdatePurchaseOrderItem{
		ID: *updatePurchaseOrderItem.ID,
	}
	jsonExisting, _ := json.Marshal(existingPurchaseOrderItem)
	json.Unmarshal(jsonExisting, &fieldsToUpdatePurchaseOrderItem.ProposedChanges)

	var updatePurchaseOrderItemMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updatePurchaseOrderItem)
	json.Unmarshal(jsonUpdate, &updatePurchaseOrderItemMap)

	updatePurchaseOrderItemTrx.mapProcessorUtility.RemoveNil(updatePurchaseOrderItemMap)

	jsonUpdate, _ = json.Marshal(updatePurchaseOrderItemMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdatePurchaseOrderItem.ProposedChanges)

	if updatePurchaseOrderItem.ProposalStatus != nil {
		fieldsToUpdatePurchaseOrderItem.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: updatePurchaseOrderItem.SubmittingAccount.ID,
		}
		if *updatePurchaseOrderItem.ProposalStatus == model.EntityProposalStatusApproved {
			json.Unmarshal(jsonUpdate, fieldsToUpdatePurchaseOrderItem)
		}
	}

	updatedPurchaseOrderItem, err := updatePurchaseOrderItemTrx.purchaseOrderItemDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdatePurchaseOrderItem.ID,
		},
		fieldsToUpdatePurchaseOrderItem,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/proposeUpdatePurchaseOrderItemComponent",
			err,
		)
	}

	return updatedPurchaseOrderItem, nil
}
