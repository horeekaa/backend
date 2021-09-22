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
	updatepurchaseOrderItemInput *model.InternalUpdatePurchaseOrderItem,
) (*model.PurchaseOrderItem, error) {
	purchaseOrderItemToUpdate := &model.DatabaseUpdatePurchaseOrderItem{}
	jsonTemp, _ := json.Marshal(updatepurchaseOrderItemInput)
	json.Unmarshal(jsonTemp, purchaseOrderItemToUpdate)

	existingPurchaseOrderItem, err := updatePurchaseOrderItemTrx.purchaseOrderItemDataSource.GetMongoDataSource().FindByID(
		purchaseOrderItemToUpdate.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updatePurchaseOrderItem",
			err,
		)
	}

	if purchaseOrderItemToUpdate.ProductVariant != nil {
		_, err := updatePurchaseOrderItemTrx.purchaseOrderItemLoader.TransactionBody(
			session,
			purchaseOrderItemToUpdate.MouItem,
			purchaseOrderItemToUpdate.ProductVariant,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updatePurchaseOrderItem",
				err,
			)
		}
		purchaseOrderItemToUpdate.UnitPrice = &purchaseOrderItemToUpdate.ProductVariant.RetailPrice
		if existingPurchaseOrderItem.MouItem != nil {
			index := funk.IndexOf(
				existingPurchaseOrderItem.MouItem.AgreedProduct.Variants,
				func(pv *model.InternalAgreedProductVariantInput) bool {
					return pv.ID == purchaseOrderItemToUpdate.ProductVariant.ID
				},
			)
			if index > -1 {
				purchaseOrderItemToUpdate.UnitPrice = &existingPurchaseOrderItem.MouItem.AgreedProduct.Variants[index].RetailPrice
			}
		}
	}
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
