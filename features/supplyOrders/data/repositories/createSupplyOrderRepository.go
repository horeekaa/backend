package supplyorderdomainrepositories

import (
	"encoding/json"

	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createSupplyOrderRepository struct {
	createSupplyOrderTransactionComponent supplyorderdomainrepositoryinterfaces.CreateSupplyOrderTransactionComponent
	createSupplyOrderItemComponent        supplyorderitemdomainrepositoryinterfaces.CreateSupplyOrderItemTransactionComponent
	mongoDBTransaction                    mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreateSupplyOrderRepository(
	createSupplyOrderRepositoryTransactionComponent supplyorderdomainrepositoryinterfaces.CreateSupplyOrderTransactionComponent,
	createSupplyOrderItemComponent supplyorderitemdomainrepositoryinterfaces.CreateSupplyOrderItemTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (supplyorderdomainrepositoryinterfaces.CreateSupplyOrderRepository, error) {
	createSupplyOrderRepo := &createSupplyOrderRepository{
		createSupplyOrderRepositoryTransactionComponent,
		createSupplyOrderItemComponent,
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
			supplyOrderItem.SupplyOrder = &model.ObjectIDOnly{
				ID: &generatedObjectID,
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
	return (output).(*model.SupplyOrder), nil
}
