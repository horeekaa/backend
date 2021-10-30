package purchaseordertosupplydomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type createPurchaseOrderToSupplyRepository struct {
	purchaseOrderDataSource   databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
	createPOToSupplyComponent purchaseordertosupplydomainrepositoryinterfaces.CreatePurchaseOrderToSupplyTransactionComponent
	mongoDBTransaction        mongodbcoretransactioninterfaces.MongoRepoTransaction
}

func NewCreatePurchaseOrderToSupplyRepository(
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
	createPOToSupplyComponent purchaseordertosupplydomainrepositoryinterfaces.CreatePurchaseOrderToSupplyTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (purchaseordertosupplydomainrepositoryinterfaces.CreatePurchaseOrderToSupplyRepository, error) {
	createPurchaseOrderToSupplyRepo := &createPurchaseOrderToSupplyRepository{
		purchaseOrderDataSource:   purchaseOrderDataSource,
		createPOToSupplyComponent: createPOToSupplyComponent,
		mongoDBTransaction:        mongoDBTransaction,
	}

	mongoDBTransaction.SetTransaction(
		createPurchaseOrderToSupplyRepo,
		"CreatePurchaseOrderToSupplyRepository",
	)

	return createPurchaseOrderToSupplyRepo, nil
}
func (createPOTSRepo *createPurchaseOrderToSupplyRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return createPOTSRepo.createPOToSupplyComponent.PreTransaction(
		input.(*model.PurchaseOrder),
	)
}

func (createPOTSRepo *createPurchaseOrderToSupplyRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	purchaseOrder := input.(*model.PurchaseOrder)

	return createPOTSRepo.createPOToSupplyComponent.TransactionBody(
		operationOption,
		purchaseOrder,
	)
}

func (createPOTSRepo *createPurchaseOrderToSupplyRepository) RunTransaction() ([]*model.PurchaseOrderToSupply, error) {
	purchaseOrders, err := createPOTSRepo.purchaseOrderDataSource.GetMongoDataSource().Find(
		map[string]interface{}{
			"status": model.PurchaseOrderStatusConfirmed,
		},
		&mongodbcoretypes.PaginationOptions{},
		nil,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createPurchaseOrderToSupplyRepository",
			err,
		)
	}

	purchaseOrderToSupplyOutput := []*model.PurchaseOrderToSupply{}

	for _, po := range purchaseOrders {
		purchaseOrderToSupplies, err := createPOTSRepo.mongoDBTransaction.RunTransaction(po)
		if err != nil {
			return nil, err
		}
		purchaseOrderToSupplyOutput = append(
			purchaseOrderToSupplyOutput,
			purchaseOrderToSupplies.([]*model.PurchaseOrderToSupply)...,
		)
	}

	return purchaseOrderToSupplyOutput, nil
}
