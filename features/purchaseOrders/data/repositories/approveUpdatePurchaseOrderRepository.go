package purchaseorderdomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type approveUpdatePurchaseOrderRepository struct {
	approveUpdatepurchaseOrderTransactionComponent purchaseorderdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderTransactionComponent
	mongoDBTransaction                             mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewApproveUpdatePurchaseOrderRepository(
	approveUpdatepurchaseOrderTransactionComponent purchaseorderdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (purchaseorderdomainrepositoryinterfaces.ApproveUpdatePurchaseOrderRepository, error) {
	approveUpdatePurchaseOrderRepo := &approveUpdatePurchaseOrderRepository{
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

	return approveUpdatePurchaseOrderRepo.approveUpdatepurchaseOrderTransactionComponent.TransactionBody(
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
