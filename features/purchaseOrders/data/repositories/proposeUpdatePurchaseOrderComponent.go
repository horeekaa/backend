package purchaseorderdomainrepositories

import (
	"encoding/json"
	"math"
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
	pathIdentity                string
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
		pathIdentity:                "ProposeUpdatePurchaseOrderComponent",
	}, nil
}

func (updatePurchaseOrderTrx *proposeUpdatePurchaseOrderTransactionComponent) PreTransaction(
	input *model.InternalUpdatePurchaseOrder,
) (*model.InternalUpdatePurchaseOrder, error) {
	return input, nil
}

func (updatePurchaseOrderTrx *proposeUpdatePurchaseOrderTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdatePurchaseOrder,
) (*model.PurchaseOrder, error) {
	updatePurchaseOrder := &model.DatabaseUpdatePurchaseOrder{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updatePurchaseOrder)

	existingPurchaseOrder, err := updatePurchaseOrderTrx.purchaseOrderDataSource.GetMongoDataSource().FindByID(
		updatePurchaseOrder.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updatePurchaseOrderTrx.pathIdentity,
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
			updatePurchaseOrderTrx.pathIdentity,
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
			updatePurchaseOrderTrx.pathIdentity,
			err,
		)
	}

	totalPrice := 0
	totalReturn := 0
	for _, item := range purchaseOrderItems {
		if !item.CustomerAgreed ||
			(item.ProposalStatus == model.EntityProposalStatusProposed &&
				item.ProposedChanges.ProposalStatus == model.EntityProposalStatusRejected) {
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

	updatePurchaseOrder.FinalSalesAmount = func(i int) *int { return &i }(totalSales)

	if existingPurchaseOrder.Mou != nil {
		mouId := *existingPurchaseOrder.Mou.ID
		if updatePurchaseOrder.Mou != nil {
			mouId = updatePurchaseOrder.Mou.ID
		}
		existingMou, err := updatePurchaseOrderTrx.mouDataSource.GetMongoDataSource().FindByID(
			mouId,
			session,
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				updatePurchaseOrderTrx.pathIdentity,
				err,
			)
		}

		if *updatePurchaseOrder.Total < *existingMou.MinimumOrderValueBeforeDelivery {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.POMinimumOrderValueHasNotMet,
				updatePurchaseOrderTrx.pathIdentity,
				nil,
			)
		}

		existingMou.RemainingCreditLimit -= *updatePurchaseOrder.FinalSalesAmount - existingPurchaseOrder.FinalSalesAmount
		if existingMou.RemainingCreditLimit < 0 {
			return nil, horeekaacorefailure.NewFailureObject(
				horeekaacorefailureenums.POSalesAmountExceedCreditLimit,
				updatePurchaseOrderTrx.pathIdentity,
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
				updatePurchaseOrderTrx.pathIdentity,
				err,
			)
		}
	}

	if input.MarkAsReceived != nil {
		if *input.MarkAsReceived {
			loc, _ := time.LoadLocation("Asia/Bangkok")
			currentTime := time.Now().UTC()
			updatePurchaseOrder.ReceivingDateTime = &currentTime
			updatePurchaseOrder.Status = func(m model.PurchaseOrderStatus) *model.PurchaseOrderStatus {
				return &m
			}(model.PurchaseOrderStatusWaitingForInvoice)
			updatePurchaseOrder.ProposalStatus = func(m model.EntityProposalStatus) *model.EntityProposalStatus {
				return &m
			}(model.EntityProposalStatusApproved)

			if existingPurchaseOrder.Type == model.PurchaseOrderTypeMouBased {
				existingMou, err := updatePurchaseOrderTrx.mouDataSource.GetMongoDataSource().FindByID(
					*existingPurchaseOrder.Mou.ID,
					session,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						updatePurchaseOrderTrx.pathIdentity,
						err,
					)
				}

				paymentDueDate := currentTime.In(loc).AddDate(
					0, 0,
					*existingMou.PaymentCompletionLimitInDays,
				)

				lastDateOfMonth := time.Date(
					paymentDueDate.Year(),
					paymentDueDate.Month()+1,
					1, 0, 0, 0, 0,
					paymentDueDate.Location(),
				).AddDate(0, 0, -1)

				paymentDay := math.Min(
					float64(paymentDueDate.Day()+15-(paymentDueDate.Day()%15)+1),
					float64(lastDateOfMonth.Day()),
				)

				paymentDueDate = time.Date(
					paymentDueDate.Year(),
					paymentDueDate.Month(),
					int(paymentDay), 0, 0, 0, 0,
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
			updatePurchaseOrderTrx.pathIdentity,
			err,
		)
	}
	updatePurchaseOrder.RecentLog = &model.ObjectIDOnly{ID: &loggingOutput.ID}

	currentTime := time.Now().UTC()
	updatePurchaseOrder.UpdatedAt = &currentTime

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
			updatePurchaseOrderTrx.pathIdentity,
			err,
		)
	}

	return updatedPurchaseOrder, nil
}
