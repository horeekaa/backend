package purchaseorderdomainrepositories

import (
	"encoding/json"
	"time"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	databasepurchaseorderItemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type updatePurchaseOrderByCronRepository struct {
	memberAccessDataSource      databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	mouDataSource               databasemoudatasourceinterfaces.MouDataSource
	purchaseOrderDataSource     databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
	purchaseOrderItemDataSource databasepurchaseorderItemdatasourceinterfaces.PurchaseOrderItemDataSource
	createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent
	pathIdentity                string
}

func NewUpdatePurchaseOrderByCronRepository(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	mouDataSource databasemoudatasourceinterfaces.MouDataSource,
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
	purchaseOrderItemDataSource databasepurchaseorderItemdatasourceinterfaces.PurchaseOrderItemDataSource,
	createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
) (purchaseorderdomainrepositoryinterfaces.UpdatePurchaseOrderByCronRepository, error) {
	updatePurchaseOrderByCronRepo := &updatePurchaseOrderByCronRepository{
		memberAccessDataSource,
		mouDataSource,
		purchaseOrderDataSource,
		purchaseOrderItemDataSource,
		createNotificationComponent,
		"UpdatePurchaseOrderByCronRepository",
	}

	return updatePurchaseOrderByCronRepo, nil
}

func (updatePurchaseOrderByCronRepo *updatePurchaseOrderByCronRepository) RunTransaction() ([]*model.PurchaseOrder, error) {
	currentDateTime := time.Now().UTC()
	futureDateOnly := time.Date(
		currentDateTime.Year(),
		currentDateTime.Month(),
		currentDateTime.Day()-3,
		0, 0, 0, 0,
		currentDateTime.Location(),
	)

	purchaseOrders, err := updatePurchaseOrderByCronRepo.purchaseOrderDataSource.GetMongoDataSource().Find(
		map[string]interface{}{
			"status": model.PurchaseOrderStatusProcessed,
			"updatedAt": map[string]interface{}{
				"$lte": futureDateOnly,
			},
		},
		&mongodbcoretypes.PaginationOptions{},
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updatePurchaseOrderByCronRepo.pathIdentity,
			err,
		)
	}

	updatedPurchaseOrders := []*model.PurchaseOrder{}
	for _, purchaseOrder := range purchaseOrders {
		purchaseOrderItems, err := updatePurchaseOrderByCronRepo.purchaseOrderItemDataSource.GetMongoDataSource().Find(
			map[string]interface{}{
				"_id": map[string]interface{}{
					"$in": funk.Map(
						purchaseOrder.Items,
						func(po *model.ObjectIDOnly) interface{} {
							return po.ID
						},
					),
				},
				"deliveryDetail.status": model.DeliveryStatusDelivered,
			},
			&mongodbcoretypes.PaginationOptions{},
			&mongodbcoretypes.OperationOptions{},
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				updatePurchaseOrderByCronRepo.pathIdentity,
				err,
			)
		}
		if len(purchaseOrderItems) != len(purchaseOrder.Items) {
			continue
		}

		loc, _ := time.LoadLocation("Asia/Bangkok")
		currentTime := time.Now().UTC()
		updatePurchaseOrder := &model.DatabaseUpdatePurchaseOrder{}
		updatePurchaseOrder.ReceivingDateTime = &currentTime
		updatePurchaseOrder.Status = func(m model.PurchaseOrderStatus) *model.PurchaseOrderStatus {
			return &m
		}(model.PurchaseOrderStatusWaitingForInvoice)

		if purchaseOrder.Type == model.PurchaseOrderTypeMouBased {
			existingMou, err := updatePurchaseOrderByCronRepo.mouDataSource.GetMongoDataSource().FindByID(
				*purchaseOrder.Mou.ID,
				&mongodbcoretypes.OperationOptions{},
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					updatePurchaseOrderByCronRepo.pathIdentity,
					err,
				)
			}

			paymentDueDate := currentTime.In(loc).AddDate(
				0, 0,
				*existingMou.PaymentCompletionLimitInDays,
			)

			paymentDay := paymentDueDate.Day() + 15 - (paymentDueDate.Day() % 15)
			paymentDueDate = time.Date(
				paymentDueDate.Year(),
				paymentDueDate.Month(),
				paymentDay, 0, 0, 0, 0,
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

		updatedPurchaseOrder, err := updatePurchaseOrderByCronRepo.purchaseOrderDataSource.GetMongoDataSource().Update(
			map[string]interface{}{
				"_id": purchaseOrder.ID,
			},
			updatePurchaseOrder,
			&mongodbcoretypes.OperationOptions{},
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				updatePurchaseOrderByCronRepo.pathIdentity,
				err,
			)
		}
		updatedPurchaseOrders = append(updatedPurchaseOrders, updatedPurchaseOrder)

		go func() {
			memberAccesses, err := updatePurchaseOrderByCronRepo.memberAccessDataSource.GetMongoDataSource().Find(
				map[string]interface{}{
					"memberAccessRefType": model.MemberAccessRefTypeOrganizationsBased,
					"status":              model.MemberAccessStatusActive,
					"proposalStatus":      model.EntityProposalStatusApproved,
					"invitationAccepted":  true,
					"organization._id":    updatedPurchaseOrder.Organization.ID,
				},
				&mongodbcoretypes.PaginationOptions{},
				&mongodbcoretypes.OperationOptions{},
			)
			if err != nil {
				return
			}

			for _, memberAccess := range memberAccesses {
				notificationToCreate := &model.InternalCreateNotification{
					NotificationCategory: model.NotificationCategoryPurchaseOrderUpdatedReceived,
					PayloadOptions: &model.PayloadOptionsInput{
						PurchaseOrderPayload: &model.PurchaseOrderPayloadInput{
							PurchaseOrder: &model.PurchaseOrderForNotifPayloadInput{},
						},
					},
					RecipientAccount: &model.ObjectIDOnly{
						ID: &memberAccess.Account.ID,
					},
				}

				jsonTemp, _ := json.Marshal(updatedPurchaseOrder)
				json.Unmarshal(jsonTemp, &notificationToCreate.PayloadOptions.PurchaseOrderPayload.PurchaseOrder)

				_, err = updatePurchaseOrderByCronRepo.createNotificationComponent.TransactionBody(
					&mongodbcoretypes.OperationOptions{},
					notificationToCreate,
				)
				if err != nil {
					return
				}
			}
		}()
	}

	return updatedPurchaseOrders, err
}
