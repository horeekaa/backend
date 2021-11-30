package purchaseorderitemdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasepurchaseOrderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	purchaseorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type proposeUpdatePurchaseOrderItemTransactionComponent struct {
	purchaseOrderItemDataSource databasepurchaseOrderitemdatasourceinterfaces.PurchaseOrderItemDataSource
	loggingDataSource           databaseloggingdatasourceinterfaces.LoggingDataSource
	purchaseOrderItemLoader     purchaseorderitemdomainrepositoryutilityinterfaces.PurchaseOrderItemLoader
	mapProcessorUtility         coreutilityinterfaces.MapProcessorUtility
}

func NewProposeUpdatePurchaseOrderItemTransactionComponent(
	purchaseOrderItemDataSource databasepurchaseOrderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	purchaseOrderItemLoader purchaseorderitemdomainrepositoryutilityinterfaces.PurchaseOrderItemLoader,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
) (purchaseorderitemdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderItemTransactionComponent, error) {
	return &proposeUpdatePurchaseOrderItemTransactionComponent{
		purchaseOrderItemDataSource: purchaseOrderItemDataSource,
		loggingDataSource:           loggingDataSource,
		purchaseOrderItemLoader:     purchaseOrderItemLoader,
		mapProcessorUtility:         mapProcessorUtility,
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

	if updatePurchaseOrderItem.DeliveryDetail == nil {
		updatePurchaseOrderItem.DeliveryDetail = &model.InternalUpdatePurchaseOrderItemDelivery{}
	}
	_, err = updatePurchaseOrderItemTrx.purchaseOrderItemLoader.TransactionBody(
		session,
		updatePurchaseOrderItem.MouItem,
		updatePurchaseOrderItem.ProductVariant,
		updatePurchaseOrderItem.DeliveryDetail.Address,
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
