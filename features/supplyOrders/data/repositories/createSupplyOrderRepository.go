package supplyorderdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databasememberaccessdatasourceinterfaces "github.com/horeekaa/backend/features/memberAccesses/data/dataSources/databases/interfaces/sources"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createSupplyOrderRepository struct {
	memberAccessDataSource                databasememberaccessdatasourceinterfaces.MemberAccessDataSource
	createSupplyOrderTransactionComponent supplyorderdomainrepositoryinterfaces.CreateSupplyOrderTransactionComponent
	createSupplyOrderItemComponent        supplyorderitemdomainrepositoryinterfaces.CreateSupplyOrderItemTransactionComponent
	createNotificationComponent           notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent
	mongoDBTransaction                    mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreateSupplyOrderRepository(
	memberAccessDataSource databasememberaccessdatasourceinterfaces.MemberAccessDataSource,
	createSupplyOrderRepositoryTransactionComponent supplyorderdomainrepositoryinterfaces.CreateSupplyOrderTransactionComponent,
	createSupplyOrderItemComponent supplyorderitemdomainrepositoryinterfaces.CreateSupplyOrderItemTransactionComponent,
	createNotificationComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (supplyorderdomainrepositoryinterfaces.CreateSupplyOrderRepository, error) {
	createSupplyOrderRepo := &createSupplyOrderRepository{
		memberAccessDataSource,
		createSupplyOrderRepositoryTransactionComponent,
		createSupplyOrderItemComponent,
		createNotificationComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		createSupplyOrderRepo,
		"CreateSupplyOrderRepository",
	)

	return createSupplyOrderRepo, nil
}

func (createSupplyOrderRepo *createSupplyOrderRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return createSupplyOrderRepo.createSupplyOrderTransactionComponent.PreTransaction(
		input.(*model.InternalCreateSupplyOrder),
	)
}

func (createSupplyOrderRepo *createSupplyOrderRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	supplyOrderToCreate := input.(*model.InternalCreateSupplyOrder)

	if len(supplyOrderToCreate.Items) > 0 {
		generatedObjectID := createSupplyOrderRepo.createSupplyOrderTransactionComponent.GenerateNewObjectID()
		savedSupplyOrderItems := []*model.InternalCreateSupplyOrderItem{}
		for _, supplyOrderItem := range supplyOrderToCreate.Items {
			if *supplyOrderToCreate.MemberAccess.Organization.Type == model.OrganizationTypePartner {
				supplyOrderItem.PartnerAgreed = func(b bool) *bool { return &b }(true)
			}
			supplyOrderItem.SupplyOrder = &model.SupplyOrderForSupplyOrderItemInput{
				ID: generatedObjectID,
				Organization: &model.ObjectIDOnly{
					ID: supplyOrderToCreate.MemberAccess.Organization.ID,
				},
			}
			supplyOrderItem.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
				return &s
			}(*supplyOrderToCreate.ProposalStatus)
			supplyOrderItem.SubmittingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
				return &m
			}(*supplyOrderToCreate.SubmittingAccount)
			createdSupplyOrderOutput, err := createSupplyOrderRepo.createSupplyOrderItemComponent.TransactionBody(
				operationOption,
				supplyOrderItem,
			)
			if err != nil {
				return nil, err
			}
			savedsupplyOrderItem := &model.InternalCreateSupplyOrderItem{}
			jsonTemp, _ := json.Marshal(createdSupplyOrderOutput)
			json.Unmarshal(jsonTemp, savedsupplyOrderItem)
			savedSupplyOrderItems = append(savedSupplyOrderItems, savedsupplyOrderItem)
		}
		supplyOrderToCreate.Items = savedSupplyOrderItems
	}

	return createSupplyOrderRepo.createSupplyOrderTransactionComponent.TransactionBody(
		operationOption,
		supplyOrderToCreate,
	)
}

func (createSupplyOrderRepo *createSupplyOrderRepository) RunTransaction(
	input *model.InternalCreateSupplyOrder,
) (*model.SupplyOrder, error) {
	output, err := createSupplyOrderRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}

	createdSupplyOrderOutput := (output).(*model.SupplyOrder)
	go func() {
		memberAccesses, err := createSupplyOrderRepo.memberAccessDataSource.GetMongoDataSource().Find(
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
				NotificationCategory: model.NotificationCategorySupplyOrderCreated,
				PayloadOptions: &model.PayloadOptionsInput{
					SupplyOrderPayload: &model.SupplyOrderPayloadInput{
						SupplyOrder: &model.SupplyOrderForNotifPayloadInput{},
					},
				},
				RecipientAccount: &model.ObjectIDOnly{
					ID: &memberAccess.Account.ID,
				},
			}

			jsonTemp, _ := json.Marshal(createdSupplyOrderOutput)
			json.Unmarshal(jsonTemp, &notificationToCreate.PayloadOptions.SupplyOrderPayload.SupplyOrder)

			_, err = createSupplyOrderRepo.createNotificationComponent.TransactionBody(
				&mongodbcoretypes.OperationOptions{},
				notificationToCreate,
			)
			if err != nil {
				return
			}
		}
	}()

	return createdSupplyOrderOutput, nil
}
