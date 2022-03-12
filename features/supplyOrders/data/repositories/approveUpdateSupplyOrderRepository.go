package supplyorderdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	databasesupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdateSupplyOrderRepository struct {
	supplyOrderDataSource                        databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource
	approveUpdatesupplyOrderItemComponent        supplyorderitemdomainrepositoryinterfaces.ApproveUpdateSupplyOrderItemTransactionComponent
	approveUpdateSupplyOrderTransactionComponent supplyorderdomainrepositoryinterfaces.ApproveUpdateSupplyOrderTransactionComponent
	mongoDBTransaction                           mongodbcoretransactioninterfaces.MongoRepoTransaction
	pathIdentity                                 string
}

func NewApproveUpdateSupplyOrderRepository(
	supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource,
	approveUpdatesupplyOrderItemComponent supplyorderitemdomainrepositoryinterfaces.ApproveUpdateSupplyOrderItemTransactionComponent,
	approveUpdateSupplyOrderTransactionComponent supplyorderdomainrepositoryinterfaces.ApproveUpdateSupplyOrderTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (supplyorderdomainrepositoryinterfaces.ApproveUpdateSupplyOrderRepository, error) {
	approveUpdateSupplyOrderRepo := &approveUpdateSupplyOrderRepository{
		supplyOrderDataSource,
		approveUpdatesupplyOrderItemComponent,
		approveUpdateSupplyOrderTransactionComponent,
		mongoDBTransaction,
		"ApproveUpdateSupplyOrderRepository",
	}

	mongoDBTransaction.SetTransaction(
		approveUpdateSupplyOrderRepo,
		"ApproveUpdateSupplyOrderRepository",
	)

	return approveUpdateSupplyOrderRepo, nil
}

func (approveUpdateSupplyOrderRepo *approveUpdateSupplyOrderRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return input, nil
}

func (approveUpdateSupplyOrderRepo *approveUpdateSupplyOrderRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	supplyOrderToApprove := input.(*model.InternalUpdateSupplyOrder)
	existingSupplyOrder, err := approveUpdateSupplyOrderRepo.supplyOrderDataSource.GetMongoDataSource().FindByID(
		supplyOrderToApprove.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			approveUpdateSupplyOrderRepo.pathIdentity,
			err,
		)
	}
	if existingSupplyOrder.ProposedChanges.ProposalStatus == model.EntityProposalStatusProposed {
		if existingSupplyOrder.ProposedChanges.Items != nil {
			for _, soItem := range existingSupplyOrder.ProposedChanges.Items {
				updateSupplyOrderItem := &model.InternalUpdateSupplyOrderItem{
					ID: &soItem.ID,
				}
				updateSupplyOrderItem.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*supplyOrderToApprove.RecentApprovingAccount)
				updateSupplyOrderItem.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*supplyOrderToApprove.ProposalStatus)

				_, err := approveUpdateSupplyOrderRepo.approveUpdatesupplyOrderItemComponent.TransactionBody(
					operationOption,
					updateSupplyOrderItem,
				)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return approveUpdateSupplyOrderRepo.approveUpdateSupplyOrderTransactionComponent.TransactionBody(
		operationOption,
		supplyOrderToApprove,
	)
}

func (approveUpdateSupplyOrderRepo *approveUpdateSupplyOrderRepository) RunTransaction(
	input *model.InternalUpdateSupplyOrder,
) (*model.SupplyOrder, error) {
	output, err := approveUpdateSupplyOrderRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).(*model.SupplyOrder), err
}
