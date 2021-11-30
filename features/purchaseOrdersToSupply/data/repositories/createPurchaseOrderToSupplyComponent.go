package purchaseordertosupplydomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	databasepurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type createPurchaseOrderToSupplyTransactionComponent struct {
	purchaseOrderDataSource         databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
	purchaseOrderItemDataSource     databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource
}

func NewCreatePurchaseOrderToSupplyTransactionComponent(
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
	purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource,
) (purchaseordertosupplydomainrepositoryinterfaces.CreatePurchaseOrderToSupplyTransactionComponent, error) {
	return &createPurchaseOrderToSupplyTransactionComponent{
		purchaseOrderDataSource:         purchaseOrderDataSource,
		purchaseOrderItemDataSource:     purchaseOrderItemDataSource,
		purchaseOrderToSupplyDataSource: purchaseOrderToSupplyDataSource,
	}, nil
}

func (createPOToSupplyTrx *createPurchaseOrderToSupplyTransactionComponent) PreTransaction(
	createPurchaseOrderInput *model.PurchaseOrder,
) (*model.PurchaseOrder, error) {
	return createPurchaseOrderInput, nil
}

func (createPOToSupplyTrx *createPurchaseOrderToSupplyTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.PurchaseOrder,
) ([]*model.PurchaseOrderToSupply, error) {
	updatedPOToSupplies := []*model.PurchaseOrderToSupply{}

	for _, poItem := range input.Items {
		existingPOItem, err := createPOToSupplyTrx.purchaseOrderItemDataSource.GetMongoDataSource().FindByID(
			poItem.ID,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/createPurchaseOrderToSupplyComponent",
				err,
			)
		}

		existingPurchaseOrderToSupply, err := createPOToSupplyTrx.purchaseOrderToSupplyDataSource.GetMongoDataSource().FindOne(
			map[string]interface{}{
				"productVariant._id":       existingPOItem.ProductVariant.ID,
				"expectedDeliverySchedule": existingPOItem.DeliveryDetail.ExpectedDeliverySchedule,
				"addressRegionGroup._id":   existingPOItem.DeliveryDetail.Address.AddressRegionGroup.ID,
				"status":                   model.PurchaseOrderToSupplyStatusCummulating,
			},
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/createPurchaseOrderToSupplyComponent",
				err,
			)
		}

		if existingPurchaseOrderToSupply == nil {
			poToCreate := &model.DatabaseCreatePurchaseOrderToSupply{
				ProductVariant:           &model.ProductVariantForPurchaseOrderItemInput{},
				AddressRegionGroup:       &model.AddressRegionGroupForPurchaseOrderToSupplyInput{},
				ExpectedDeliverySchedule: *existingPOItem.DeliveryDetail.ExpectedDeliverySchedule,
				Status:                   func(s model.PurchaseOrderToSupplyStatus) *model.PurchaseOrderToSupplyStatus { return &s }(model.PurchaseOrderToSupplyStatusCummulating),
			}

			jsonTemp, _ := json.Marshal(existingPOItem.ProductVariant)
			json.Unmarshal(jsonTemp, &poToCreate.ProductVariant)

			jsonTemp, _ = json.Marshal(existingPOItem.DeliveryDetail.Address.AddressRegionGroup)
			json.Unmarshal(jsonTemp, &poToCreate.AddressRegionGroup)

			jsonTemp, _ = json.Marshal(
				map[string]interface{}{
					"Tags": funk.Map(
						existingPOItem.ProductVariant.Product.Taggings,
						func(t *model.TaggingForPurchaseOrderItem) interface{} {
							return t.Tag
						},
					),
				},
			)
			json.Unmarshal(jsonTemp, poToCreate)

			existingPurchaseOrderToSupply, err = createPOToSupplyTrx.purchaseOrderToSupplyDataSource.GetMongoDataSource().Create(
				poToCreate,
				session,
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					"/createPurchaseOrderToSupplyComponent",
					err,
				)
			}
		}

		updatedPOToSupply, err := createPOToSupplyTrx.purchaseOrderToSupplyDataSource.GetMongoDataSource().Update(
			map[string]interface{}{
				"_id": existingPurchaseOrderToSupply.ID,
			},
			&model.DatabaseUpdatePurchaseOrderToSupply{
				QuantityRequested: func(i int) *int { return &i }(existingPurchaseOrderToSupply.QuantityRequested + existingPOItem.Quantity),
				PurchaseOrderItems: funk.Map(
					append(
						existingPurchaseOrderToSupply.PurchaseOrderItems,
						&model.PurchaseOrderItem{
							ID: poItem.ID,
						},
					),
					func(m *model.PurchaseOrderItem) *model.ObjectIDOnly {
						return &model.ObjectIDOnly{
							ID: &m.ID,
						}
					},
				).([]*model.ObjectIDOnly),
			},
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/createPurchaseOrderToSupplyComponent",
				err,
			)
		}

		_, err = createPOToSupplyTrx.purchaseOrderItemDataSource.GetMongoDataSource().Update(
			map[string]interface{}{
				"_id": existingPOItem.ID,
			},
			&model.DatabaseUpdatePurchaseOrderItem{
				Status: func(m model.PurchaseOrderItemStatus) *model.PurchaseOrderItemStatus {
					return &m
				}(model.PurchaseOrderItemStatusAwaitingFulfillment),
				PurchaseOrderToSupply: &model.ObjectIDOnly{
					ID: &updatedPOToSupply.ID,
				},
			},
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/createPurchaseOrderToSupplyComponent",
				err,
			)
		}

		updatedPOToSupplies = append(updatedPOToSupplies, updatedPOToSupply)
	}
	_, err := createPOToSupplyTrx.purchaseOrderDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": input.ID,
		},
		&model.DatabaseUpdatePurchaseOrder{
			Status: func(m model.PurchaseOrderStatus) *model.PurchaseOrderStatus {
				return &m
			}(model.PurchaseOrderStatusProcessed),
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createPurchaseOrderToSupplyComponent",
			err,
		)
	}

	return updatedPOToSupplies, nil
}
