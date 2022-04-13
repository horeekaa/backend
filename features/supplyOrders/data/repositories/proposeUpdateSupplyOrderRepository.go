package supplyorderdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	databasesupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	"github.com/horeekaa/backend/model"
	"github.com/thoas/go-funk"
)

type proposeUpdateSupplyOrderRepository struct {
	memberAccessDataSource                       databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	supplyOrderItemDataSource                    databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource
	supplyOrderDataSource                        databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource
	proposeUpdateSupplyOrderTransactionComponent supplyorderdomainrepositoryinterfaces.ProposeUpdateSupplyOrderTransactionComponent
	createSupplyOrderItemComponent               supplyorderitemdomainrepositoryinterfaces.CreateSupplyOrderItemTransactionComponent
	proposeUpdateSupplyOrderItemComponent        supplyorderitemdomainrepositoryinterfaces.ProposeUpdateSupplyOrderItemTransactionComponent
	approveUpdateSupplyOrderItemComponent        supplyorderitemdomainrepositoryinterfaces.ApproveUpdateSupplyOrderItemTransactionComponent
	createNotificationComponent                  notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent
	mongoDBTransaction                           mongodbcoretransactioninterfaces.MongoRepoTransaction
	pathIdentity                                 string
}

func NewProposeUpdateSupplyOrderRepository(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource,
	supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource,
	proposeUpdateSupplyOrderRepositoryTransactionComponent supplyorderdomainrepositoryinterfaces.ProposeUpdateSupplyOrderTransactionComponent,
	createSupplyOrderItemComponent supplyorderitemdomainrepositoryinterfaces.CreateSupplyOrderItemTransactionComponent,
	proposeUpdateSupplyOrderItemComponent supplyorderitemdomainrepositoryinterfaces.ProposeUpdateSupplyOrderItemTransactionComponent,
	approveUpdateSupplyOrderItemComponent supplyorderitemdomainrepositoryinterfaces.ApproveUpdateSupplyOrderItemTransactionComponent,
	createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (supplyorderdomainrepositoryinterfaces.ProposeUpdateSupplyOrderRepository, error) {
	proposeUpdateSupplyOrderRepo := &proposeUpdateSupplyOrderRepository{
		memberAccessDataSource,
		supplyOrderItemDataSource,
		supplyOrderDataSource,
		proposeUpdateSupplyOrderRepositoryTransactionComponent,
		createSupplyOrderItemComponent,
		proposeUpdateSupplyOrderItemComponent,
		approveUpdateSupplyOrderItemComponent,
		createNotificationComponent,
		mongoDBTransaction,
		"ProposeUpdateSupplyOrderRepository",
	}

	mongoDBTransaction.SetTransaction(
		proposeUpdateSupplyOrderRepo,
		"ProposeUpdateSupplyOrderRepository",
	)

	return proposeUpdateSupplyOrderRepo, nil
}

func (updateSupplyOrderRepo *proposeUpdateSupplyOrderRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return updateSupplyOrderRepo.proposeUpdateSupplyOrderTransactionComponent.PreTransaction(
		input.(*model.InternalUpdateSupplyOrder),
	)
}

func (updateSupplyOrderRepo *proposeUpdateSupplyOrderRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	supplyOrderToUpdate := input.(*model.InternalUpdateSupplyOrder)
	existingsupplyOrder, err := updateSupplyOrderRepo.supplyOrderDataSource.GetMongoDataSource().FindByID(
		supplyOrderToUpdate.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateSupplyOrderRepo.pathIdentity,
			err,
		)
	}

	if supplyOrderToUpdate.Items != nil {
		savedsupplyOrderItems := existingsupplyOrder.Items
		for _, supplyOrderItemToUpdate := range supplyOrderToUpdate.Items {
			if supplyOrderItemToUpdate.ID != nil {
				if !funk.Contains(
					existingsupplyOrder.Items,
					func(mi *model.SupplyOrderItem) bool {
						return mi.ID == *supplyOrderItemToUpdate.ID
					},
				) {
					continue
				}
				if supplyOrderItemToUpdate.ProposalStatus != nil {
					supplyOrderItemToUpdate.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
						return &m
					}(*supplyOrderToUpdate.SubmittingAccount)

					_, err := updateSupplyOrderRepo.approveUpdateSupplyOrderItemComponent.TransactionBody(
						operationOption,
						supplyOrderItemToUpdate,
					)
					if err != nil {
						return nil, err
					}
					continue
				}
				supplyOrderItemToUpdate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*supplyOrderToUpdate.ProposalStatus)
				supplyOrderItemToUpdate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*supplyOrderToUpdate.SubmittingAccount)
				if *supplyOrderToUpdate.MemberAccess.Organization.Type != model.OrganizationTypePartner {
					supplyOrderItemToUpdate.PartnerAgreed = func(b bool) *bool { return &b }(false)
				}

				_, err := updateSupplyOrderRepo.proposeUpdateSupplyOrderItemComponent.TransactionBody(
					operationOption,
					supplyOrderItemToUpdate,
				)
				if err != nil {
					return nil, err
				}
				continue
			}

			supplyOrderItemToCreate := &model.InternalCreateSupplyOrderItem{}
			jsonTemp, _ := json.Marshal(supplyOrderItemToUpdate)
			json.Unmarshal(jsonTemp, supplyOrderItemToCreate)
			supplyOrderItemToCreate.SupplyOrder = &model.ObjectIDOnly{
				ID: &existingsupplyOrder.ID,
			}
			supplyOrderItemToCreate.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*supplyOrderToUpdate.ProposalStatus)
			supplyOrderItemToCreate.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*supplyOrderToUpdate.SubmittingAccount)
			if *supplyOrderToUpdate.MemberAccess.Organization.Type == model.OrganizationTypePartner {
				supplyOrderItemToCreate.PartnerAgreed = func(b bool) *bool { return &b }(true)
			}
			for i, descPhoto := range supplyOrderItemToUpdate.Photos {
				supplyOrderItemToCreate.Photos[i].Photo.File = descPhoto.Photo.File
			}

			savedSupplyOrderItem, err := updateSupplyOrderRepo.createSupplyOrderItemComponent.TransactionBody(
				operationOption,
				supplyOrderItemToCreate,
			)
			if err != nil {
				return nil, err
			}
			savedsupplyOrderItems = append(savedsupplyOrderItems, savedSupplyOrderItem)
		}
		jsonTemp, _ := json.Marshal(
			map[string]interface{}{
				"Items": savedsupplyOrderItems,
			},
		)
		json.Unmarshal(jsonTemp, supplyOrderToUpdate)
	}

	return updateSupplyOrderRepo.proposeUpdateSupplyOrderTransactionComponent.TransactionBody(
		operationOption,
		supplyOrderToUpdate,
	)
}

func (updateSupplyOrderRepo *proposeUpdateSupplyOrderRepository) RunTransaction(
	input *model.InternalUpdateSupplyOrder,
) (*model.SupplyOrder, error) {
	output, err := updateSupplyOrderRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}

	updatedSupplyOrderOutput := (output).(*model.SupplyOrder)

	existingSupplyOrder, err := updateSupplyOrderRepo.supplyOrderDataSource.GetMongoDataSource().FindByID(
		input.ID,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			updateSupplyOrderRepo.pathIdentity,
			err,
		)
	}

	for _, supplyOrderItemToUpdate := range input.Items {
		if supplyOrderItemToUpdate.ID != nil {
			if !funk.Contains(
				existingSupplyOrder.Items,
				func(mi *model.SupplyOrderItem) bool {
					return mi.ID == *supplyOrderItemToUpdate.ID
				},
			) {
				continue
			}

			existingSupplyOrderItem, err := updateSupplyOrderRepo.supplyOrderItemDataSource.GetMongoDataSource().FindByID(
				*supplyOrderItemToUpdate.ID,
				&mongodbcoretypes.OperationOptions{},
			)
			if err != nil {
				return nil, horeekaacoreexceptiontofailure.ConvertException(
					updateSupplyOrderRepo.pathIdentity,
					err,
				)
			}

			memberAccessQuery := map[string]interface{}{}
			var notifCategory model.NotificationCategory
			if supplyOrderItemToUpdate.ProposalStatus != nil {
				memberAccessQuery = map[string]interface{}{
					"memberAccessRefType": model.MemberAccessRefTypeOrganizationsBased,
					"status":              model.MemberAccessStatusActive,
					"proposalStatus":      model.EntityProposalStatusApproved,
					"invitationAccepted":  true,
					"organization._id":    existingSupplyOrder.Organization.ID,
				}
				notifCategory = model.NotificationCategorySupplyOrderItemApproval
			}

			if supplyOrderItemToUpdate.PartnerAgreed != nil {
				memberAccessQuery = map[string]interface{}{
					"memberAccessRefType": model.MemberAccessRefTypeOrganizationsBased,
					"status":              model.MemberAccessStatusActive,
					"proposalStatus":      model.EntityProposalStatusApproved,
					"invitationAccepted":  true,
					"organization.type":   model.OrganizationTypeInternal,
				}
				notifCategory = model.NotificationCategorySupplyOrderItemPartnerAgreed
			}

			if supplyOrderItemToUpdate.QuantityAccepted != nil {
				memberAccessQuery = map[string]interface{}{
					"memberAccessRefType": model.MemberAccessRefTypeOrganizationsBased,
					"status":              model.MemberAccessStatusActive,
					"proposalStatus":      model.EntityProposalStatusApproved,
					"invitationAccepted":  true,
					"organization._id":    existingSupplyOrder.Organization.ID,
				}
				notifCategory = model.NotificationCategorySupplyOrderItemAccepted
			}

			go func() {
				memberAccesses, err := updateSupplyOrderRepo.memberAccessDataSource.GetMongoDataSource().Find(
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
							SupplyOrderItemPayload: &model.SupplyOrderItemPayloadInput{
								SupplyOrderItem: &model.SupplyOrderItemForNotifPayloadInput{
									SupplyOrder: &model.SupplyOrderForNotifPayloadInput{},
								},
							},
						},
						RecipientAccount: &model.ObjectIDOnly{
							ID: &memberAccess.Account.ID,
						},
					}

					jsonPOItem, _ := json.Marshal(existingSupplyOrderItem)
					json.Unmarshal(jsonPOItem, &notificationToCreate.PayloadOptions.SupplyOrderItemPayload.SupplyOrderItem)

					jsonPO, _ := json.Marshal(existingSupplyOrder)
					json.Unmarshal(jsonPO, &notificationToCreate.PayloadOptions.SupplyOrderItemPayload.SupplyOrderItem.SupplyOrder)

					_, err = updateSupplyOrderRepo.createNotificationComponent.TransactionBody(
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

		existingSupplyOrderItem, err := updateSupplyOrderRepo.supplyOrderItemDataSource.GetMongoDataSource().FindOne(
			map[string]interface{}{
				"purchaseOrderToSupply._id": supplyOrderItemToUpdate.PurchaseOrderToSupply.ID,
				"supplyOrder._id":           existingSupplyOrder.ID,
			},
			&mongodbcoretypes.OperationOptions{},
		)
		if err != nil {
			return nil, horeekaacoreexceptiontofailure.ConvertException(
				updateSupplyOrderRepo.pathIdentity,
				err,
			)
		}
		if existingSupplyOrderItem == nil {
			continue
		}

		go func() {
			memberAccesses, err := updateSupplyOrderRepo.memberAccessDataSource.GetMongoDataSource().Find(
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
					NotificationCategory: model.NotificationCategorySupplyOrderItemCreated,
					PayloadOptions: &model.PayloadOptionsInput{
						SupplyOrderItemPayload: &model.SupplyOrderItemPayloadInput{
							SupplyOrderItem: &model.SupplyOrderItemForNotifPayloadInput{
								SupplyOrder: &model.SupplyOrderForNotifPayloadInput{},
							},
						},
					},
					RecipientAccount: &model.ObjectIDOnly{
						ID: &memberAccess.Account.ID,
					},
				}

				jsonPOItem, _ := json.Marshal(existingSupplyOrderItem)
				json.Unmarshal(jsonPOItem, &notificationToCreate.PayloadOptions.SupplyOrderItemPayload.SupplyOrderItem)

				jsonPO, _ := json.Marshal(existingSupplyOrder)
				json.Unmarshal(jsonPO, &notificationToCreate.PayloadOptions.SupplyOrderItemPayload.SupplyOrderItem.SupplyOrder)

				_, err = updateSupplyOrderRepo.createNotificationComponent.TransactionBody(
					&mongodbcoretypes.OperationOptions{},
					notificationToCreate,
				)
				if err != nil {
					return
				}
			}
		}()
	}

	return updatedSupplyOrderOutput, nil
}
