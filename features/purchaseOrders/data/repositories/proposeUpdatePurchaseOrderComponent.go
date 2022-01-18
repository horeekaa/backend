package purchaseorderdomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	purchaseorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type proposeUpdatePurchaseOrderTransactionComponent struct {
	purchaseOrderDataSource     databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
	purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource
	loggingDataSource           databaseloggingdatasourceinterfaces.LoggingDataSource
	mouDataSource               databasemoudatasourceinterfaces.MouDataSource
	mapProcessorUtility         coreutilityinterfaces.MapProcessorUtility
	purchaseOrderDataLoader     purchaseorderdomainrepositoryutilityinterfaces.PurchaseOrderLoader
}

func NewProposeUpdatePurchaseOrderTransactionComponent(
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
	purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mouDataSource databasemoudatasourceinterfaces.MouDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
	purchaseOrderDataLoader purchaseorderdomainrepositoryutilityinterfaces.PurchaseOrderLoader,
) (purchaseorderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderTransactionComponent, error) {
	return &proposeUpdatePurchaseOrderTransactionComponent{
		purchaseOrderDataSource:     purchaseOrderDataSource,
		purchaseOrderItemDataSource: purchaseOrderItemDataSource,
		loggingDataSource:           loggingDataSource,
		mouDataSource:               mouDataSource,
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
		nil,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/updatePurchaseOrder",
			err,
		)
	}

	purchaseOrderItems, err := updatePurchaseOrderTrx.purchaseOrderItemDataSource.GetMongoDataSource().Find(
		map[string]interface{}{
			"purchaseOrder._id": existingPurchaseOrder.ID,
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
	totalReturn := 0
	for _, item := range purchaseOrderItems {
		if item.ProposalStatus == model.EntityProposalStatusRejected || !item.CustomerAgreed {
			continue
		}
		totalPrice += item.SubTotal
		if item.PurchaseOrderItemReturn != nil {
			totalReturn += item.PurchaseOrderItemReturn.SubTotal
		}
	}
	updatePurchaseOrder.Total = &totalPrice
	updatePurchaseOrder.TotalReturn = &totalReturn
	totalSales := totalPrice - totalReturn

	totalDiscounted := existingPurchaseOrder.TotalDiscounted
	if updatePurchaseOrder.TotalDiscounted != nil {
		totalDiscounted = *updatePurchaseOrder.TotalDiscounted
	}

	discountInPercent := existingPurchaseOrder.DiscountInPercent
	if updatePurchaseOrder.DiscountInPercent != nil {
		discountInPercent = *updatePurchaseOrder.DiscountInPercent
	}

	if discountInPercent > 0 {
		totalDiscounted = totalSales * discountInPercent
	}

	updatePurchaseOrder.FinalSalesAmount = func(i int) *int { return &i }(totalSales - totalDiscounted)

	if existingPurchaseOrder.Mou != nil {
		mouId := existingPurchaseOrder.Mou.ID
		if updatePurchaseOrder.Mou != nil {
			mouId = updatePurchaseOrder.Mou.ID
		}
		existingMou, err := updatePurchaseOrderTrx.mouDataSource.GetMongoDataSource().FindByID(
			mouId,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updatePurchaseOrder",
				err,
			)
		}

		existingMou.RemainingCreditLimit -= *updatePurchaseOrder.FinalSalesAmount - existingPurchaseOrder.FinalSalesAmount
		if existingMou.RemainingCreditLimit < 0 {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.POSalesAmountExceedCreditLimit,
				"/updatePurchaseOrder",
				nil,
			)
		}

		_, err = updatePurchaseOrderTrx.mouDataSource.GetMongoDataSource().Update(
			map[string]interface{}{
				"_id": mouId,
			},
			&model.DatabaseUpdateMou{
				RemainingCreditLimit: &existingMou.RemainingCreditLimit,
			},
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				"/updatePurchaseOrder",
				err,
			)
		}
	}

	if updatePurchaseOrder.MarkAsReceived != nil {
		if *updatePurchaseOrder.MarkAsReceived {
			loc, _ := time.LoadLocation("Asia/Bangkok")
			currentTime := time.Now().UTC()
			updatePurchaseOrder.ReceivingDateTime = &currentTime
			updatePurchaseOrder.Status = func(m model.PurchaseOrderStatus) *model.PurchaseOrderStatus {
				return &m
			}(model.PurchaseOrderStatusPaymentNeeded)

			if existingPurchaseOrder.Type == model.PurchaseOrderTypeMouBased {
				existingMou, err := updatePurchaseOrderTrx.mouDataSource.GetMongoDataSource().FindByID(
					existingPurchaseOrder.Mou.ID,
					session,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/updatePurchaseOrder",
						err,
					)
				}

				paymentDueDate := currentTime.In(loc).AddDate(
					0, 0,
					*existingMou.PaymentCompletionLimitInDays,
				)
				paymentDueDate = time.Date(
					paymentDueDate.Year(),
					paymentDueDate.Month()+1,
					1, 0, 0, 0, 0,
					paymentDueDate.Location(),
				).UTC()

				updatePurchaseOrder.PaymentDueDate = &paymentDueDate
			} else {
				currentTime = currentTime.In(loc)
				currentTime = time.Date(
					currentTime.Year(),
					currentTime.Month(),
					currentTime.Day(), 0, 0, 0, 0,
					currentTime.Location(),
				).UTC()

				updatePurchaseOrder.PaymentDueDate = &currentTime
			}
		}
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
