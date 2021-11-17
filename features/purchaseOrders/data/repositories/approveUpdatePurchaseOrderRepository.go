package purchaseorderdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdatePurchaseOrderRepository struct {
	purchaseOrderDataSource                        databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
	approveUpdatePurchaseOrderItemComponent        purchaseorderitemdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderItemTransactionComponent
	approveUpdatePurchaseOrderTransactionComponent purchaseorderdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderTransactionComponent
	mongoDBTransaction                             mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewApproveUpdatePurchaseOrderRepository(
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
	approveUpdatePurchaseOrderItemComponent purchaseorderitemdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderItemTransactionComponent,
	approveUpdatepurchaseOrderTransactionComponent purchaseorderdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (purchaseorderdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderRepository, error) {
	approveUpdatePurchaseOrderRepo := &approveUpdatePurchaseOrderRepository{
		purchaseOrderDataSource,
		approveUpdatePurchaseOrderItemComponent,
		approveUpdatepurchaseOrderTransactionComponent,
		mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		approveUpdatePurchaseOrderRepo,
		"ApproveUpdatePurchaseOrderRepository",
	)

	return approveUpdatePurchaseOrderRepo, nil
}

func (approveUpdatePurchaseOrderRepo *approveUpdatePurchaseOrderRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return input, nil
}

func (approveUpdatePurchaseOrderRepo *approveUpdatePurchaseOrderRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	purchaseOrderToApprove := input.(*model.InternalUpdatePurchaseOrder)
	existingPurchaseOrder, err := approveUpdatePurchaseOrderRepo.purchaseOrderDataSource.GetMongoDataSource().FindByID(
		purchaseOrderToApprove.ID,
		operationOption,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/approveUpdatePurchaseOrderRepository",
			err,
		)
	}
	if existingPurchaseOrder.ProposedChanges.ProposalStatus == model.EntityProposalStatusProposed {
		if existingPurchaseOrder.ProposedChanges.Items != nil {
			for _, poItem := range existingPurchaseOrder.ProposedChanges.Items {
				updatePurchaseOrderItem := &model.InternalUpdatePurchaseOrderItem{
					ID: &poItem.ID,
				}
				updatePurchaseOrderItem.RecentApprovingAccount = func(m model.ObjectIDOnly) *model.ObjectIDOnly {
					return &m
				}(*purchaseOrderToApprove.RecentApprovingAccount)
				updatePurchaseOrderItem.ProposalStatus = func(s model.EntityProposalStatus) *model.EntityProposalStatus {
					return &s
				}(*purchaseOrderToApprove.ProposalStatus)

				_, err := approveUpdatePurchaseOrderRepo.approveUpdatePurchaseOrderItemComponent.TransactionBody(
					operationOption,
					updatePurchaseOrderItem,
				)
				if err != nil {
					return nil, horeekaacoreexceptiontofailure.ConvertException(
						"/approveUpdatePurchaseOrderRepository",
						err,
					)
				}
			}
		}
	}

	return approveUpdatePurchaseOrderRepo.approveUpdatePurchaseOrderTransactionComponent.TransactionBody(
		operationOption,
		purchaseOrderToApprove,
	)
}

func (approveUpdatePurchaseOrderRepo *approveUpdatePurchaseOrderRepository) RunTransaction(
	input *model.InternalUpdatePurchaseOrder,
) (*model.PurchaseOrder, error) {
	output, err := approveUpdatePurchaseOrderRepo.mongoDBTransaction.RunTransaction(input)
	if err != nil {
		return nil, err
	}
	return (output).(*model.PurchaseOrder), err
}
