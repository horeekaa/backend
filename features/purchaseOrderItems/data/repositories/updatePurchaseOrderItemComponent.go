package purchaseorderitemdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasepurchaseOrderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	purchaseorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type updatePurchaseOrderItemTransactionComponent struct {
	purchaseOrderItemDataSource databasepurchaseOrderitemdatasourceinterfaces.PurchaseOrderItemDataSource
	purchaseOrderItemLoader     purchaseorderitemdomainrepositoryutilityinterfaces.PurchaseOrderItemLoader
}

func NewUpdatePurchaseOrderItemTransactionComponent(
	purchaseOrderItemDataSource databasepurchaseOrderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
	purchaseOrderItemLoader purchaseorderitemdomainrepositoryutilityinterfaces.PurchaseOrderItemLoader,
) (purchaseorderitemdomainrepositoryinterfaces.UpdatePurchaseOrderItemTransactionComponent, error) {
	return &updatePurchaseOrderItemTransactionComponent{
		purchaseOrderItemDataSource: purchaseOrderItemDataSource,
		purchaseOrderItemLoader:     purchaseOrderItemLoader,
	}, nil
}

func (updatePurchaseOrderItemTrx *updatePurchaseOrderItemTransactionComponent) PreTransaction(
	input *model.InternalUpdatePurchaseOrderItem,
) (*model.InternalUpdatePurchaseOrderItem, error) {
	return input, nil
}

func (updatePurchaseOrderItemTrx *updatePurchaseOrderItemTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updatePurchaseOrderItemInput *model.InternalUpdatePurchaseOrderItem,
) (*model.PurchaseOrderItem, error) {
	existingPurchaseOrderItem, err := updatePurchaseOrderItemTrx.purchaseOrderItemDataSource.GetMongoDataSource().FindByID(
		*updatePurchaseOrderItemInput.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updatePurchaseOrderItem",
			err,
		)
	}

	if updatePurchaseOrderItemInput.ProductVariant != nil {
		_, err := updatePurchaseOrderItemTrx.purchaseOrderItemLoader.TransactionBody(
			session,
			updatePurchaseOrderItemInput.MouItem,
			updatePurchaseOrderItemInput.ProductVariant,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updatePurchaseOrderItem",
				err,
			)
		}
		updatePurchaseOrderItemInput.UnitPrice = &updatePurchaseOrderItemInput.ProductVariant.RetailPrice
		if existingPurchaseOrderItem.MouItem != nil {
			index := funk.IndexOf(
				existingPurchaseOrderItem.MouItem.AgreedProduct.Variants,
				func(pv *model.InternalAgreedProductVariantInput) bool {
					return pv.ID == updatePurchaseOrderItemInput.ProductVariant.ID
				},
			)
			if index > -1 {
				updatePurchaseOrderItemInput.UnitPrice = &existingPurchaseOrderItem.MouItem.AgreedProduct.Variants[index].RetailPrice
			}
		}
	}
	purchaseOrderItemToUpdate := &model.DatabaseUpdatePurchaseOrderItem{}
	jsonTemp, _ := json.Marshal(updatePurchaseOrderItemInput)
	json.Unmarshal(jsonTemp, purchaseOrderItemToUpdate)

	unitPrice := existingPurchaseOrderItem.UnitPrice
	if purchaseOrderItemToUpdate.UnitPrice != nil {
		unitPrice = *purchaseOrderItemToUpdate.UnitPrice
	}

	quantity := existingPurchaseOrderItem.Quantity
	if purchaseOrderItemToUpdate.Quantity != nil {
		quantity = *purchaseOrderItemToUpdate.Quantity
	}
	subTotal := quantity * unitPrice

	purchaseOrderItemToUpdate.SubTotal = &subTotal

	updatedPurchaseOrderItem, err := updatePurchaseOrderItemTrx.purchaseOrderItemDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": purchaseOrderItemToUpdate.ID,
		},
		purchaseOrderItemToUpdate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updatePurchaseOrderItem",
			err,
		)
	}

	return updatedPurchaseOrderItem, nil
}
