package purchaseordertosupplydomainrepositories

import (
	mongodbcoretransactioninterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/transaction"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	notificationdomainrepositoryinterfaces "github.com/horeekaa/backend/features/notifications/domain/repositories"
	databasepurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
	"github.com/horeekaa/backend/model"
)

type processPurchaseOrderToSupplyRepository struct {
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource
	processPOToSupplyComponent      purchaseordertosupplydomainrepositoryinterfaces.ProcessPurchaseOrderToSupplyTransactionComponent
	createNotifComponent            notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent
	mongoDBTransaction              mongodbcoretransactioninterfaces.MongoRepoTransaction
	pathIdentity                    string
}

func NewProcessPurchaseOrderToSupplyRepository(
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource,
	processPOToSupplyComponent purchaseordertosupplydomainrepositoryinterfaces.ProcessPurchaseOrderToSupplyTransactionComponent,
	createNotifComponent notificationdomainrepositoryinterfaces.CreateNotificationTransactionComponent,
	mongoDBTransaction mongodbcoretransactioninterfaces.MongoRepoTransaction,
) (purchaseordertosupplydomainrepositoryinterfaces.ProcessPurchaseOrderToSupplyRepository, error) {
	processPurchaseOrderToSupplyRepo := &processPurchaseOrderToSupplyRepository{
		purchaseOrderToSupplyDataSource: purchaseOrderToSupplyDataSource,
		processPOToSupplyComponent:      processPOToSupplyComponent,
		createNotifComponent:            createNotifComponent,
		mongoDBTransaction:              mongoDBTransaction,
		pathIdentity:                    "ProcessPurchaseOrderToSupplyRepository",
	}

	mongoDBTransaction.SetTransaction(
		processPurchaseOrderToSupplyRepo,
		"ProcessPurchaseOrderToSupplyRepository",
	)

	return processPurchaseOrderToSupplyRepo, nil
}

func (processPOTSRepo *processPurchaseOrderToSupplyRepository) PreTransaction(
	input interface{},
) (interface{}, error) {
	return input, nil
}

func (processPOTSRepo *processPurchaseOrderToSupplyRepository) TransactionBody(
	operationOption *mongodbcoretypes.OperationOptions,
	input interface{},
) (interface{}, error) {
	poToSupplyInput := input.(*model.PurchaseOrderToSupply)

	notifsToCreate, err := processPOTSRepo.processPOToSupplyComponent.TransactionBody(
		operationOption,
		poToSupplyInput,
	)
	if err != nil {
		return nil, err
	}

	for _, notifToCreate := range notifsToCreate {
		_, err := processPOTSRepo.createNotifComponent.TransactionBody(
			operationOption,
			notifToCreate,
		)
		if err != nil {
			return nil, err
		}
	}

	return true, nil
}

func (processPOTSRepo *processPurchaseOrderToSupplyRepository) RunTransaction() (bool, error) {
	purchaseOrdersToSupply, err := processPOTSRepo.purchaseOrderToSupplyDataSource.GetMongoDataSource().Find(
		map[string]interface{}{
			"status": model.PurchaseOrderToSupplyStatusCummulating,
		},
		&mongodbcoretypes.PaginationOptions{},
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return false, horeekaacoreexceptiontofailure.ConvertException(
			processPOTSRepo.pathIdentity,
			err,
		)
	}

	for _, poToSupply := range purchaseOrdersToSupply {
		_, err := processPOTSRepo.mongoDBTransaction.RunTransaction(poToSupply)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
