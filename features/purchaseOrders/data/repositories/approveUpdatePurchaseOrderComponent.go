package purchaseorderdomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	coreutilityinterfaces "github.com/horeekaa/backend/core/utilities/interfaces"
	databaseloggingdatasourceinterfaces "github.com/horeekaa/backend/features/loggings/data/dataSources/databases/interfaces"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	purchaseorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type approveUpdatePurchaseOrderTransactionComponent struct {
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
	loggingDataSource       databaseloggingdatasourceinterfaces.LoggingDataSource
	mapProcessorUtility     coreutilityinterfaces.MapProcessorUtility
	mouDataSource           databasemoudatasourceinterfaces.MouDataSource
	purchaseOrderDataLoader purchaseorderdomainrepositoryutilityinterfaces.PurchaseOrderLoader
	pathIdentity            string
}

func NewApproveUpdatePurchaseOrderTransactionComponent(
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
	loggingDataSource databaseloggingdatasourceinterfaces.LoggingDataSource,
	mapProcessorUtility coreutilityinterfaces.MapProcessorUtility,
	mouDataSource databasemoudatasourceinterfaces.MouDataSource,
	purchaseOrderDataLoader purchaseorderdomainrepositoryutilityinterfaces.PurchaseOrderLoader,
) (purchaseorderdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderTransactionComponent, error) {
	return &approveUpdatePurchaseOrderTransactionComponent{
		purchaseOrderDataSource: purchaseOrderDataSource,
		loggingDataSource:       loggingDataSource,
		mapProcessorUtility:     mapProcessorUtility,
		mouDataSource:           mouDataSource,
		purchaseOrderDataLoader: purchaseOrderDataLoader,
		pathIdentity:            "ApproveUpdatePurchaseOrderComponent",
	}, nil
}

func (approvePurchaseOrderTrx *approveUpdatePurchaseOrderTransactionComponent) PreTransaction(
	input *model.InternalUpdatePurchaseOrder,
) (*model.InternalUpdatePurchaseOrder, error) {
	return input, nil
}

func (approvePurchaseOrderTrx *approveUpdatePurchaseOrderTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	input *model.InternalUpdatePurchaseOrder,
) (*model.PurchaseOrder, error) {
	updatePurchaseOrder := &model.DatabaseUpdatePurchaseOrder{}
	jsonTemp, _ := json.Marshal(input)
	json.Unmarshal(jsonTemp, updatePurchaseOrder)

	existingPurchaseOrder, err := approvePurchaseOrderTrx.purchaseOrderDataSource.GetMongoDataSource().FindByID(
		updatePurchaseOrder.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approvePurchaseOrderTrx.pathIdentity,
			err,
		)
	}

	_, err = approvePurchaseOrderTrx.purchaseOrderDataLoader.TransactionBody(
		session,
		updatePurchaseOrder.Mou,
		nil,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approvePurchaseOrderTrx.pathIdentity,
			err,
		)
	}

	previousLog, err := approvePurchaseOrderTrx.loggingDataSource.GetMongoDataSource().FindByID(
		existingPurchaseOrder.RecentLog.ID,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approvePurchaseOrderTrx.pathIdentity,
			err,
		)
	}

	logToCreate := &model.CreateLogging{
		Collection: previousLog.Collection,
		CreatedByAccount: &model.ObjectIDOnly{
			ID: updatePurchaseOrder.RecentApprovingAccount.ID,
		},
		Activity:       previousLog.Activity,
		ProposalStatus: *updatePurchaseOrder.ProposalStatus,
	}
	jsonTemp, _ = json.Marshal(
		map[string]interface{}{
			"NewDocumentJSON": previousLog.NewDocumentJSON,
			"OldDocumentJSON": previousLog.OldDocumentJSON,
		},
	)
	json.Unmarshal(jsonTemp, logToCreate)

	createdLog, err := approvePurchaseOrderTrx.loggingDataSource.GetMongoDataSource().Create(
		logToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approvePurchaseOrderTrx.pathIdentity,
			err,
		)
	}

	updatePurchaseOrder.RecentLog = &model.ObjectIDOnly{ID: &createdLog.ID}

	currentTime := time.Now().UTC()
	updatePurchaseOrder.UpdatedAt = &currentTime

	fieldsToUpdatePurchaseOrder := &model.DatabaseUpdatePurchaseOrder{
		ID: updatePurchaseOrder.ID,
	}
	jsonExisting, _ := json.Marshal(existingPurchaseOrder.ProposedChanges)
	json.Unmarshal(jsonExisting, &fieldsToUpdatePurchaseOrder.ProposedChanges)

	var updatePurchaseOrderMap map[string]interface{}
	jsonUpdate, _ := json.Marshal(updatePurchaseOrder)
	json.Unmarshal(jsonUpdate, &updatePurchaseOrderMap)

	approvePurchaseOrderTrx.mapProcessorUtility.RemoveNil(updatePurchaseOrderMap)

	jsonUpdate, _ = json.Marshal(updatePurchaseOrderMap)
	json.Unmarshal(jsonUpdate, &fieldsToUpdatePurchaseOrder.ProposedChanges)

	if updatePurchaseOrder.ProposalStatus != nil {
		if *updatePurchaseOrder.ProposalStatus == model.EntityProposalStatusApproved {
			fieldsToUpdatePurchaseOrder.ProposedChanges.Status = func(m model.PurchaseOrderStatus) *model.PurchaseOrderStatus {
				return &m
			}(model.PurchaseOrderStatusProcessed)
			jsonUpdate, _ := json.Marshal(fieldsToUpdatePurchaseOrder.ProposedChanges)
			json.Unmarshal(jsonUpdate, fieldsToUpdatePurchaseOrder)
		}
		if *updatePurchaseOrder.ProposalStatus == model.EntityProposalStatusRejected {
			if existingPurchaseOrder.Mou != nil {
				mouId := *existingPurchaseOrder.Mou.ID
				if updatePurchaseOrder.Mou != nil {
					mouId = updatePurchaseOrder.Mou.ID
				}
				existingMou, err := approvePurchaseOrderTrx.mouDataSource.GetMongoDataSource().FindByID(
					mouId,
					session,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						approvePurchaseOrderTrx.pathIdentity,
						err,
					)
				}
				existingMou.RemainingCreditLimit += existingPurchaseOrder.FinalSalesAmount
				_, err = approvePurchaseOrderTrx.mouDataSource.GetMongoDataSource().Update(
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
						approvePurchaseOrderTrx.pathIdentity,
						err,
					)
				}
			}
		}
	}

	updatedPurchaseOrder, err := approvePurchaseOrderTrx.purchaseOrderDataSource.GetMongoDataSource().Update(
		map[string]interface{}{
			"_id": fieldsToUpdatePurchaseOrder.ID,
		},
		fieldsToUpdatePurchaseOrder,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approvePurchaseOrderTrx.pathIdentity,
			err,
		)
	}

	return updatedPurchaseOrder, nil
}
