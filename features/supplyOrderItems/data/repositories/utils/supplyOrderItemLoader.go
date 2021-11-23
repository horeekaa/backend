package supplyorderitemdomainrepositoryutilities

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databasepurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type supplyOrderItemLoader struct {
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource
}

func NewSupplyOrderItemLoader(
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource,
) (supplyorderitemdomainrepositoryutilityinterfaces.SupplyOrderItemLoader, error) {
	return &supplyOrderItemLoader{
		purchaseOrderToSupplyDataSource,
	}, nil
}

func (supOrderItemLoader *supplyOrderItemLoader) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	purchaseOrderToSupply *model.PurchaseOrderToSupplyForSupplyOrderItemInput,
) (bool, error) {
	if purchaseOrderToSupply != nil {
		loadedPOToSupply, err := supOrderItemLoader.purchaseOrderToSupplyDataSource.GetMongoDataSource().FindByID(
			purchaseOrderToSupply.ID,
			session,
		)
		if err != nil {
			return false, err
		}

		jsonTemp, _ := json.Marshal(loadedPOToSupply)
		json.Unmarshal(jsonTemp, purchaseOrderToSupply)
	}

	return true, nil
}
