package purchaseorderdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	purchaseorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type proposeUpdatePurchaseOrderTransactionComponent struct {
	purchaseOrderDataSource     databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
	purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource
	loggingDataSource           databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility         coreutilityinterfaces.MapProcessorUtility
	purchaseOrderDataLoader     purchaseorderdomainrepositoryutilityinterfaces.PurchaseOrderLoader
}

func NewProposeUpdatePurchaseOrderTransactionComponent(
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
	purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
	purchaseOrderDataLoader purchaseorderdomainrepositoryutilityinterfaces.PurchaseOrderLoader,
) (purchaseorderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderTransactionComponent, error) {
	return &proposeUpdatePurchaseOrderTransactionComponent{
		purchaseOrderDataSource:     purchaseOrderDataSource,
		purchaseOrderItemDataSource: purchaseOrderItemDataSource,
		loggingDataSource:           loggingDataSource,
		mapProcessorUtility:         mapProcessorUtility,
		purchaseOrderDataLoader:     purchaseOrderDataLoader,
	}, nil
}

func (updatePurchaseOrderTrx *proposeUpdatePurchaseOrderTransactionComponent) PreTransaction(
	input *model.InternalUpdatePurchaseOrder,
) (*model.InternalUpdatePurchaseOrder, error) {
	return input, nil
}

func (updatePurchaseOrderTrx *proposeUpdatePurchaseOrderTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	updatePurchaseOrder *model.InternalUpdatePurchaseOrder,
) (*model.PurchaseOrder, error) {
	existingPurchaseOrder, err := updatePurchaseOrderTrx.purchaseOrderDataSource.GetMongoDataSource().FindByID(
		updatePurchaseOrder.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updatePurchaseOrder",
			err,
		)
	}

	_, err = updatePurchaseOrderTrx.purchaseOrderDataLoader.TransactionBody(
		session,
		updatePurchaseOrder.Mou,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updatePurchaseOrder",
			err,
		)
	}

	if updatePurchaseOrder.Items != nil {
		purchaseOrderItems, err := updatePurchaseOrderTrx.purchaseOrderItemDataSource.GetMongoDataSource().Find(
			map[string]interface{}{
				"_id": map[string]interface{}{
					"$in": funk.Map(
						updatePurchaseOrder.Items,
						func(it *model.InternalUpdatePurchaseOrderItem) interface{} {
							return it.ID
						},
					),
				},
			},
			&mongodbcoretypes.PaginationOptions{},
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updatePurchaseOrder",
				err,
			)
		}

		totalPrice := 0
		for _, item := range purchaseOrderItems {
			totalPrice += item.SubTotal
		}
		updatePurchaseOrder.Total = &totalPrice

		totalDiscounted := existingPurchaseOrder.TotalDiscounted
		if updatePurchaseOrder.TotalDiscounted != nil {
			totalDiscounted = *updatePurchaseOrder.TotalDiscounted
		}

		discountInPercent := existingPurchaseOrder.DiscountInPercent
		if updatePurchaseOrder.DiscountInPercent != nil {
			discountInPercent = *updatePurchaseOrder.DiscountInPercent
		}

		if discountInPercent > 0 {
			totalDiscounted = totalPrice * discountInPercent
		}

		updatePurchaseOrder.FinalSalesAmount = func(i int) *int { return &i }(totalPrice - totalDiscounted)
	}

	newDocumentJson, _ := json.Marshal(*updatePurchaseOrder)
	oldDocumentJson, _ := json.Marshal(*existingPurchaseOrder)
	loggingOutput, err := updatePurchaseOrderTrx.loggingDataSource.GetMongoDataSource().Create(
		&model.CreateLogging{
			Collection: "PurchaseOrder",
			Document: &model.ObjectIDOnly{
				ID: &existingPurchaseOrder.ID,
			},
			NewDocumentJSON: func(s string) *string { return &s }(string(newDocumentJson)),
			OldDocumentJSON: func(s string) *string { return &s }(string(oldDocumentJson)),
			CreatedByAccount: &model.ObjectIDOnly{
				ID: updatePurchaseOrder.SubmittingAccount.ID,
			},
			Activity:       model.LoggedActivityUpdate,
			ProposalStatus: *updatePurchaseOrder.ProposalStatus,
		},
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updatePurchaseOrder",
			err,
		)
	}
	updatePurchaseOrder.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	fieldsToUpdatePurchaseOrder := &model.DatabaseUpdatePurchaseOrder{
		ID: updatePurchaseOrder.ID,
	}
	jsonExisting, _ := json.Marshal(existingPurchaseOrder)
	json.Unmarshal(jsonExisting, &fieldsToUpdatePurchaseOrder.ProposedChanges)

	var updatepurchaseOrderMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updatePurchaseOrder)
	json.Unmarshal(jsonUpdate, &updatepurchaseOrderMap)

	updatePurchaseOrderTrx.mapProcessorUtility.RemoveNil(updatepurchaseOrderMap)

	jsonUpdate, _ = json.Marshal(updatepurchaseOrderMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdatePurchaseOrder.ProposedChanges)
	jsonUpdate, err = json.Marshal(fieldsToUpdatePurchaseOrder.ProposedChanges)

	if updatePurchaseOrder.ProposalStatus != nil {
		fieldsToUpdatePurchaseOrder.RecentApprovingAccount = &model.ObjectIDOnly{
			ID: updatePurchaseOrder.SubmittingAccount.ID,
		}
		if *updatePurchaseOrder.ProposalStatus == model.EntityProposalStatusApproved {
			json.Unmarshal(jsonUpdate, fieldsToUpdatePurchaseOrder)
		}
	}

	updatedPurchaseOrder, err := updatePurchaseOrderTrx.purchaseOrderDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdatePurchaseOrder.ID,
		},
		fieldsToUpdatePurchaseOrder,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updatePurchaseOrder",
			err,
		)
	}

	return updatedPurchaseOrder, nil
}
