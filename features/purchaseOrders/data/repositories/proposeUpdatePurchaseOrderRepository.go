package purchaseorderdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacorefailure "github.com/horeekaa/backend/core/errors/failures"
	horeekaacorefailureenums "github.com/horeekaa/backend/core/errors/failures/enums"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	invoicedomainrepositoryinterfaces "github.com/horeekaa/backend/features/invoices/domain/repositories"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	databasepurchaseorderItemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type proposeUpdatePurchaseOrderRepository struct {
	memberAccessDataSource                         databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	purchaseOrderDataSource                        databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
	purchaseOrderItemDataSource                    databasepurchaseorderItemdatasourceinterfaces.PurchaseOrderItemDataSource
	proposeUpdatePurchaseOrderTransactionComponent purchaseorderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderTransactionComponent
	createPurchaseOrderItemComponent               purchaseorderitemdomainrepositoryinterfaces.CreatePurchaseOrderItemTransactionComponent
	updatePurchaseOrderItemComponent               purchaseorderitemdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderItemTransactionComponent
	approvePurchaseOrderItemComponent              purchaseorderitemdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderItemTransactionComponent
	updateInvoiceTrxComponent                      invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent
	createNotificationComponent                    notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent
	mongoDBTransaction                             mongodbcoretransactioninterfaces.MongoRepoTransaction
	pathIdentity                                   string
}

func NewProposeUpdatePurchaseOrderRepository(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
	purchaseOrderItemDataSource databasepurchaseorderItemdatasourceinterfaces.PurchaseOrderItemDataSource,
	proposeUpdatePurchaseOrderRepositoryTransactionComponent purchaseorderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderTransactionComponent,
	createPurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.CreatePurchaseOrderItemTransactionComponent,
	updatePurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderItemTransactionComponent,
	approvePurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderItemTransactionComponent,
	updateInvoiceTrxComponent invoicedomainrepositoryinterfaces.UpdateInvoiceTransactionComponent,
	createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (purchaseorderdomainrepositoryinterfaces.ProposeUpdatePurchaseOrderRepository, error) {
	proposeUpdatePurchaseOrderRepo := &proposeUpdatePurchaseOrderRepository{
		memberAccessDataSource,
		purchaseOrderDataSource,
		purchaseOrderItemDataSource,
		proposeUpdatePurchaseOrderRepositoryTransactionComponent,
		createPurchaseOrderItemComponent,
		updatePurchaseOrderItemComponent,
		approvePurchaseOrderItemComponent,
		updateInvoiceTrxComponent,
		createNotificationComponent,
		mongoDBTransaction,
		"ProposeUpdatePurchaseOrderRepository",
	}

	mongoDBTransaction.SetTransaction(
		proposeUpdatePurchaseOrderRepo,
		"ProposeUpdatePurchaseOrderRepository",
	)

	return proposeUpdatePurchaseOrderRepo, nil
}

func (updatePurchaseOrderRepo *proposeUpdatePurchaseOrderRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updatePurchaseOrderRepo.proposeUpdatePurchaseOrderTransactionComponent.PreTransaction(
		input.(*model.InternalUpdatePurchaseOrder),
	)
}

func (updatePurchaseOrderRepo *proposeUpdatePurchaseOrderRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	purchaseOrderToUpdate := input.(*model.InternalUpdatePurchaseOrder)
	existingPurchaseOrder, err := updatePurchaseOrderRepo.purchaseOrderDataSource.GetMongoDataSource().FindByID(
		purchaseOrderToUpdate.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updatePurchaseOrderRepo.pathIdentity,
			err,
		)
	}

	if purchaseOrderToUpdate.Items != nil {
		savedPurchaseOrderItems := existingPurchaseOrder.Items
		for _, purchaseOrderItemToUpdate := range purchaseOrderToUpdate.Items {
			if purchaseOrderItemToUpdate.ID != nil {
				if !funk.Contains(
					existingPurchaseOrder.Items,
					func(mi *model.PurchaseOrderItem) bool {
						return mi.ID == *purchaseOrderItemToUpdate.ID
					},
				) {
					continue
				}
				if purchaseOrderItemToUpdate.ProposalStatus != nil {
					purchaseOrderItemToUpdate.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
						return &m
					}(*purchaseOrderToUpdate.SubmittingAccount)

					_, err := updatePurchaseOrderRepo.approvePurchaseOrderItemComponent.TransactionBody(
						operationOption,
						purchaseOrderItemToUpdate,
					)
					if err != nil {
						return nil, err
					}
					continue
				}

				purchaseOrderItemToUpdate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*purchaseOrderToUpdate.ProposalStatus)
				purchaseOrderItemToUpdate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*purchaseOrderToUpdate.SubmittingAccount)
				if *purchaseOrderToUpdate.MemberAccess.Organization.Type != model.OrganizationTypeCustomer {
					purchaseOrderItemToUpdate.CustomerAgreed = func(b bool) *bool { return &b }(false)
				}

				_, err := updatePurchaseOrderRepo.updatePurchaseOrderItemComponent.TransactionBody(
					operationOption,
					purchaseOrderItemToUpdate,
				)
				if err != nil {
					return nil, err
				}
				continue
			}
			if existingPurchaseOrder.Mou != nil && purchaseOrderItemToUpdate.MouItem == nil {
				return nil, horeekaacorefailure.NewFailureObject(
					horeekaacorefailureenums.POItemMismatchWithPOType,
					updatePurchaseOrderRepo.pathIdentity,
					nil,
				)
			}

			purchaseOrderItemToCreate := &model.InternalCreatePurchaseOrderItem{}
			jsonTemp, _ := json.Marshal(purchaseOrderItemToUpdate)
			json.Unmarshal(jsonTemp, purchaseOrderItemToCreate)
			purchaseOrderItemToCreate.PurchaseOrder = &model.ObjectIDOnly{
				ID: &existingPurchaseOrder.ID,
			}
			purchaseOrderItemToCreate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*purchaseOrderToUpdate.ProposalStatus)
			purchaseOrderItemToCreate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*purchaseOrderToUpdate.SubmittingAccount)
			if *purchaseOrderToUpdate.MemberAccess.Organization.Type == model.OrganizationTypeCustomer {
				purchaseOrderItemToCreate.CustomerAgreed = func(b bool) *bool { return &b }(true)
			}

			savedPurchaseOrderItem, err := updatePurchaseOrderRepo.createPurchaseOrderItemComponent.TransactionBody(
				operationOption,
				purchaseOrderItemToCreate,
			)
			if err != nil {
				return nil, err
			}
			savedPurchaseOrderItems = append(savedPurchaseOrderItems, savedPurchaseOrderItem)
		}
		jsonTemp, _ := json.Marshal(
			map[string]interface{}{
				"Items": savedPurchaseOrderItems,
			},
		)
		json.Unmarshal(jsonTemp, purchaseOrderToUpdate)
	}

	if purchaseOrderToUpdate.ProposalStatus != nil {
		if *purchaseOrderToUpdate.ProposalStatus == model.EntityProposalStatusApproved {
			if purchaseOrderToUpdate.Invoice != nil {
				if funk.Get(existingPurchaseOrder, "Invoice.ID") != nil {
					if existingPurchaseOrder.Invoice.ID.Hex() != purchaseOrderToUpdate.Invoice.ID.Hex() {
						_, err := updatePurchaseOrderRepo.updateInvoiceTrxComponent.TransactionBody(
							operationOption,
							&model.InternalUpdateInvoice{
								ID: *existingPurchaseOrder.Invoice.ID,
								PurchaseOrdersToRemove: []*model.ObjectIDOnly{
									{ID: &existingPurchaseOrder.ID},
								},
							},
						)
						if err != nil {
							return nil, err
						}
					}
				}
				_, err := updatePurchaseOrderRepo.updateInvoiceTrxComponent.TransactionBody(
					operationOption,
					&model.InternalUpdateInvoice{
						ID: *purchaseOrderToUpdate.Invoice.ID,
						PurchaseOrdersToAdd: []*model.ObjectIDOnly{
							{ID: &existingPurchaseOrder.ID},
						},
					},
				)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return updatePurchaseOrderRepo.proposeUpdatePurchaseOrderTransactionComponent.TransactionBody(
		operationOption,
		purchaseOrderToUpdate,
	)
}

func (updatePurchaseOrderRepo *proposeUpdatePurchaseOrderRepository) RunTransaction(
	input *model.InternalUpdatePurchaseOrder,
) (*model.PurchaseOrder, error) {
	output, err := updatePurchaseOrderRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	updatedPurchaseOrderOutput := (output).(*model.PurchaseOrder)

	existingPurchaseOrder, err := updatePurchaseOrderRepo.purchaseOrderDataSource.GetMongoDataSource().FindByID(
		input.ID,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updatePurchaseOrderRepo.pathIdentity,
			err,
		)
	}

	for _, purchaseOrderItemToUpdate := range input.Items {
		if purchaseOrderItemToUpdate.ID != nil {
			if !funk.Contains(
				existingPurchaseOrder.Items,
				func(mi *model.PurchaseOrderItem) bool {
					return mi.ID == *purchaseOrderItemToUpdate.ID
				},
			) {
				continue
			}

			existingPurchaseOrderItem, err := updatePurchaseOrderRepo.purchaseOrderItemDataSource.GetMongoDataSource().FindByID(
				*purchaseOrderItemToUpdate.ID,
				&mongodbcoretypes.OperationOptions{},
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					updatePurchaseOrderRepo.pathIdentity,
					err,
				)
			}

			memberAccessQuery := map[string]interface{}{}
			var notifCategory model.NotificationCategory
			if purchaseOrderItemToUpdate.ProposalStatus != nil {
				memberAccessQuery = map[string]interface{}{
					"memberAccessRefType": model.MemberAccessRefTypeOrganizationsBased,
					"status":              model.MemberAccessStatusActive,
					"proposalStatus":      model.EntityProposalStatusApproved,
					"invitationAccepted":  true,
					"organization._id":    existingPurchaseOrder.Organization.ID,
				}
				notifCategory = model.NotificationCategoryPurchaseOrderItemApproval
			}

			if purchaseOrderItemToUpdate.CustomerAgreed != nil {
				memberAccessQuery = map[string]interface{}{
					"memberAccessRefType": model.MemberAccessRefTypeOrganizationsBased,
					"status":              model.MemberAccessStatusActive,
					"proposalStatus":      model.EntityProposalStatusApproved,
					"invitationAccepted":  true,
					"organization.type":   model.OrganizationTypeInternal,
				}
				notifCategory = model.NotificationCategoryPurchaseOrderItemCustomerAgreed
			}

			if purchaseOrderItemToUpdate.QuantityFulfilled != nil {
				memberAccessQuery = map[string]interface{}{
					"memberAccessRefType": model.MemberAccessRefTypeOrganizationsBased,
					"status":              model.MemberAccessStatusActive,
					"proposalStatus":      model.EntityProposalStatusApproved,
					"invitationAccepted":  true,
					"organization._id":    existingPurchaseOrder.Organization.ID,
				}
				notifCategory = model.NotificationCategoryPurchaseOrderItemFulfilled
			}

			go func() {
				memberAccesses, err := updatePurchaseOrderRepo.memberAccessDataSource.GetMongoDataSource().Find(
					memberAccessQuery,
					&mongodbcoretypes.PaginationOptions{},
					&mongodbcoretypes.OperationOptions{},
				)
				if err != nil {
					return
				}

				for _, memberAccess := range memberAccesses {
					notificationToCreate := &model.InternalCreateNotification{
						NotificationCategory: notifCategory,
						PayloadOptions: &model.PayloadOptionsInput{
							PurchaseOrderItemPayload: &model.PurchaseOrderItemPayloadInput{
								PurchaseOrderItem: &model.PurchaseOrderItemForNotifPayloadInput{
									PurchaseOrder: &model.PurchaseOrderForNotifPayloadInput{},
								},
							},
						},
						RecipientAccount: &model.ObjectIDOnly{
							ID: &memberAccess.Account.ID,
						},
					}

					jsonPOItem, _ := json.Marshal(existingPurchaseOrderItem)
					json.Unmarshal(jsonPOItem, &notificationToCreate.PayloadOptions.PurchaseOrderItemPayload.PurchaseOrderItem)

					jsonPO, _ := json.Marshal(existingPurchaseOrder)
					json.Unmarshal(jsonPO, &notificationToCreate.PayloadOptions.PurchaseOrderItemPayload.PurchaseOrderItem.PurchaseOrder)

					_, err = updatePurchaseOrderRepo.createNotificationComponent.TransactionBody(
						&mongodbcoretypes.OperationOptions{},
						notificationToCreate,
					)
					if err != nil {
						return
					}
				}
			}()
			continue
		}

		existingPurchaseOrderItem, err := updatePurchaseOrderRepo.purchaseOrderItemDataSource.GetMongoDataSource().FindOne(
			map[string]interface{}{
				"productVariant._id": purchaseOrderItemToUpdate.ProductVariant.ID,
				"quantity":           purchaseOrderItemToUpdate.Quantity,
			},
			&mongodbcoretypes.OperationOptions{},
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				updatePurchaseOrderRepo.pathIdentity,
				err,
			)
		}

		go func() {
			memberAccesses, err := updatePurchaseOrderRepo.memberAccessDataSource.GetMongoDataSource().Find(
				map[string]interface{}{
					"memberAccessRefType": model.MemberAccessRefTypeOrganizationsBased,
					"status":              model.MemberAccessStatusActive,
					"proposalStatus":      model.EntityProposalStatusApproved,
					"invitationAccepted":  true,
					"organization.type":   model.OrganizationTypeInternal,
				},
				&mongodbcoretypes.PaginationOptions{},
				&mongodbcoretypes.OperationOptions{},
			)
			if err != nil {
				return
			}

			for _, memberAccess := range memberAccesses {
				notificationToCreate := &model.InternalCreateNotification{
					NotificationCategory: model.NotificationCategoryPurchaseOrderItemCreated,
					PayloadOptions: &model.PayloadOptionsInput{
						PurchaseOrderItemPayload: &model.PurchaseOrderItemPayloadInput{
							PurchaseOrderItem: &model.PurchaseOrderItemForNotifPayloadInput{
								PurchaseOrder: &model.PurchaseOrderForNotifPayloadInput{},
							},
						},
					},
					RecipientAccount: &model.ObjectIDOnly{
						ID: &memberAccess.Account.ID,
					},
				}

				jsonPOItem, _ := json.Marshal(existingPurchaseOrderItem)
				json.Unmarshal(jsonPOItem, &notificationToCreate.PayloadOptions.PurchaseOrderItemPayload.PurchaseOrderItem)

				jsonPO, _ := json.Marshal(existingPurchaseOrder)
				json.Unmarshal(jsonPO, &notificationToCreate.PayloadOptions.PurchaseOrderItemPayload.PurchaseOrderItem.PurchaseOrder)

				_, err = updatePurchaseOrderRepo.createNotificationComponent.TransactionBody(
					&mongodbcoretypes.OperationOptions{},
					notificationToCreate,
				)
				if err != nil {
					return
				}
			}
		}()
	}

	return updatedPurchaseOrderOutput, nil
}
