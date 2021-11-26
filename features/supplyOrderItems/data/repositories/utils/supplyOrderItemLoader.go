package supplyorderitemdomainrepositoryutilities

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	databasepurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type supplyOrderItemLoader struct {
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource
	addressDataSource               databaseaddressdatasourceinterfaces.AddressDataSource
}

func NewSupplyOrderItemLoader(
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource,
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
) (supplyorderitemdomainrepositoryutilityinterfaces.SupplyOrderItemLoader, error) {
	return &supplyOrderItemLoader{
		purchaseOrderToSupplyDataSource,
		addressDataSource,
	}, nil
}

func (supOrderItemLoader *supplyOrderItemLoader) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	purchaseOrderToSupply *model.PurchaseOrderToSupplyForSupplyOrderItemInput,
	pickUpAddress *model.AddressForSupplyOrderItemInput,
) (bool, error) {
	purchaseOrderToSupplyLoadedChan := make(chan bool)
	addressLoadedChan := make(chan bool)
	errChan := make(chan error)

	go func() {
		if purchaseOrderToSupply == nil {
			purchaseOrderToSupplyLoadedChan <- true
			return
		}

		loadedPOToSupply, err := supOrderItemLoader.purchaseOrderToSupplyDataSource.GetMongoDataSource().FindByID(
			purchaseOrderToSupply.ID,
			session,
		)
		if err != nil {
			errChan <- err
			return
		}

		jsonTemp, _ := json.Marshal(loadedPOToSupply)
		json.Unmarshal(jsonTemp, purchaseOrderToSupply)

		purchaseOrderToSupplyLoadedChan <- true
	}()

	go func() {
		if pickUpAddress == nil {
			addressLoadedChan <- true
			return
		}

		loadedAddress, err := supOrderItemLoader.addressDataSource.GetMongoDataSource().FindByID(
			pickUpAddress.ID,
			session,
		)
		if err != nil {
			errChan <- err
			return
		}
		jsonTemp, _ := json.Marshal(loadedAddress)
		json.Unmarshal(jsonTemp, pickUpAddress)

		addressLoadedChan <- true
	}()

	for i := 0; i < 2; {
		select {
		case err := <-errChan:
			return false, err
		case _ = <-purchaseOrderToSupplyLoadedChan:
			i++
		case _ = <-addressLoadedChan:
			i++
		}
	}

	return true, nil
}
