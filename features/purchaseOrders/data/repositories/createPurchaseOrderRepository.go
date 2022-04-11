package purchaseorderdomainrepositories

import (
	"encoding/json"
	"reflect"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createPurchaseOrderRepository struct {
	memberAccessDataSource                  databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	createPurchaseOrderTransactionComponent purchaseorderdomainrepositoryinterfaces.CreatePurchaseOrderTransactionComponent
	createPurchaseOrderItemComponent        purchaseorderitemdomainrepositoryinterfaces.CreatePurchaseOrderItemTransactionComponent
	createNotificationComponent             notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent
	mongoDBTransaction                      mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreatePurchaseOrderRepository(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	createPurchaseOrderRepositoryTransactionComponent purchaseorderdomainrepositoryinterfaces.CreatePurchaseOrderTransactionComponent,
	createPurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.CreatePurchaseOrderItemTransactionComponent,
	createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (purchaseorderdomainrepositoryinterfaces.CreatePurchaseOrderRepository, error) {
	createPurchaseOrderRepo := &createPurchaseOrderRepository{
		memberAccessDataSource,
		createPurchaseOrderRepositoryTransactionComponent,
		createPurchaseOrderItemComponent,
		createNotificationComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		createPurchaseOrderRepo,
		"CreatePurchaseOrderRepository",
	)

	return createPurchaseOrderRepo, nil
}

func (createPurchaseOrderRepo *createPurchaseOrderRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return createPurchaseOrderRepo.createPurchaseOrderTransactionComponent.PreTransaction(
		input.(*model.InternalCreatePurchaseOrder),
	)
}

func (createPurchaseOrderRepo *createPurchaseOrderRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	purchaseOrderToCreate := input.(*model.InternalCreatePurchaseOrder)
	purchaseOrders := []*model.PurchaseOrder{}

	if len(purchaseOrderToCreate.MouItems) > 0 {
		generatedObjectID := createPurchaseOrderRepo.createPurchaseOrderTransactionComponent.GenerateNewObjectID()
		savedPurchaseOrderItems := map[string][]*model.InternalCreatePurchaseOrderItem{}
		mouForPurchaseOrders := map[string]*model.ObjectIDOnlyOutput{}
		for _, purchaseOrderItem := range purchaseOrderToCreate.MouItems {
			if purchaseOrderItem.MouItem == nil {
				continue
			}
			if *purchaseOrderToCreate.MemberAccess.Organization.Type == model.OrganizationTypeCustomer {
				purchaseOrderItem.CustomerAgreed = func(b bool) *bool { return &b }(true)
			}
			purchaseOrderItem.PurchaseOrder = &model.ObjectIDOnly{
				ID: &generatedObjectID,
			}
			purchaseOrderItem.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*purchaseOrderToCreate.ProposalStatus)
			purchaseOrderItem.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*purchaseOrderToCreate.SubmittingAccount)
			createdPurchaseOrderItemOutput, err := createPurchaseOrderRepo.createPurchaseOrderItemComponent.TransactionBody(
				operationOption,
				purchaseOrderItem,
			)
			if err != nil {
				return nil, err
			}
			savedPurchaseOrderItem := &model.InternalCreatePurchaseOrderItem{}
			jsonTemp, _ := json.Marshal(createdPurchaseOrderItemOutput)
			json.Unmarshal(jsonTemp, savedPurchaseOrderItem)

			mouStringID := createdPurchaseOrderItemOutput.MouItem.Mou.ID.Hex()
			mouForPurchaseOrders[mouStringID] = createdPurchaseOrderItemOutput.MouItem.Mou

			if savedPurchaseOrderItems[mouStringID] == nil {
				savedPurchaseOrderItems[mouStringID] = []*model.InternalCreatePurchaseOrderItem{}
			}
			savedPurchaseOrderItems[mouStringID] = append(
				savedPurchaseOrderItems[mouStringID],
				savedPurchaseOrderItem,
			)
		}

		for _, e := range reflect.ValueOf(savedPurchaseOrderItems).MapKeys() {
			purchaseOrderToCreate.Items = savedPurchaseOrderItems[e.String()]
			purchaseOrderToCreate.Type = model.PurchaseOrderTypeMouBased
			purchaseOrderToCreate.Mou = &model.ObjectIDOnly{
				ID: mouForPurchaseOrders[e.String()].ID,
			}

			purchaseOrder, err := createPurchaseOrderRepo.createPurchaseOrderTransactionComponent.TransactionBody(
				operationOption,
				purchaseOrderToCreate,
			)
			if err != nil {
				return nil, err
			}

			purchaseOrders = append(purchaseOrders, purchaseOrder)
		}
	}

	if len(purchaseOrderToCreate.RetailItems) > 0 {
		generatedObjectID := createPurchaseOrderRepo.createPurchaseOrderTransactionComponent.GenerateNewObjectID()
		savedPurchaseOrderItems := []*model.InternalCreatePurchaseOrderItem{}
		for _, purchaseOrderItem := range purchaseOrderToCreate.RetailItems {
			if *purchaseOrderToCreate.MemberAccess.Organization.Type == model.OrganizationTypeCustomer {
				purchaseOrderItem.CustomerAgreed = func(b bool) *bool { return &b }(true)
			}
			purchaseOrderItem.PurchaseOrder = &model.ObjectIDOnly{
				ID: &generatedObjectID,
			}
			purchaseOrderItem.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*purchaseOrderToCreate.ProposalStatus)
			purchaseOrderItem.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*purchaseOrderToCreate.SubmittingAccount)
			createdPurchaseOrderOutput, err := createPurchaseOrderRepo.createPurchaseOrderItemComponent.TransactionBody(
				operationOption,
				purchaseOrderItem,
			)
			if err != nil {
				return nil, err
			}
			savedPurchaseOrderItem := &model.InternalCreatePurchaseOrderItem{}
			jsonTemp, _ := json.Marshal(createdPurchaseOrderOutput)
			json.Unmarshal(jsonTemp, savedPurchaseOrderItem)
			savedPurchaseOrderItems = append(savedPurchaseOrderItems, savedPurchaseOrderItem)
		}
		purchaseOrderToCreate.Items = savedPurchaseOrderItems
		purchaseOrderToCreate.Type = model.PurchaseOrderTypeRetail
		purchaseOrder, err := createPurchaseOrderRepo.createPurchaseOrderTransactionComponent.TransactionBody(
			operationOption,
			purchaseOrderToCreate,
		)
		if err != nil {
			return nil, err
		}

		purchaseOrders = append(purchaseOrders, purchaseOrder)
	}

	return purchaseOrders, nil
}

func (createPurchaseOrderRepo *createPurchaseOrderRepository) RunTransaction(
	input *model.InternalCreatePurchaseOrder,
) ([]*model.PurchaseOrder, error) {
	output, err := createPurchaseOrderRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}

	createdPurchaseOrdersOutput := (output).([]*model.PurchaseOrder)
	go func() {
		memberAccesses, err := createPurchaseOrderRepo.memberAccessDataSource.GetMongoDataSource().Find(
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

		for _, createdPurchaseOrder := range createdPurchaseOrdersOutput {
			for _, memberAccess := range memberAccesses {
				notificationToCreate := &model.InternalCreateNotification{
					NotificationCategory: model.NotificationCategoryPurchaseOrderCreated,
					PayloadOptions: &model.PayloadOptionsInput{
						PurchaseOrderPayload: &model.PurchaseOrderPayloadInput{
							PurchaseOrder: &model.PurchaseOrderForNotifPayloadInput{},
						},
					},
					RecipientAccount: &model.ObjectIDOnly{
						ID: &memberAccess.Account.ID,
					},
				}

				jsonTemp, _ := json.Marshal(createdPurchaseOrder)
				json.Unmarshal(jsonTemp, &notificationToCreate.PayloadOptions.PurchaseOrderPayload.PurchaseOrder)

				_, err = createPurchaseOrderRepo.createNotificationComponent.TransactionBody(
					&mongodbcoretypes.OperationOptions{},
					notificationToCreate,
				)
				if err != nil {
					return
				}
			}
		}
	}()

	return createdPurchaseOrdersOutput, nil
}
