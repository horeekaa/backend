package supplyorderitemdomainrepositoryutilities

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databaseaccountdatasourceinterfaces "github.com/horeekaa/backend/features/accounts/data/dataSources/databases/interfaces/sources"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	databasepurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type supplyOrderItemLoader struct {
	accountDataSource               databaseaccountdatasourceinterfaces.AccountDataSource
	personDataSource                databaseaccountdatasourceinterfaces.PersonDataSource
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource
	addressDataSource               databaseaddressdatasourceinterfaces.AddressDataSource
}

func NewSupplyOrderItemLoader(
	accountDataSource databaseaccountdatasourceinterfaces.AccountDataSource,
	personDataSource databaseaccountdatasourceinterfaces.PersonDataSource,
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource,
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
) (supplyorderitemdomainrepositoryutilityinterfaces.SupplyOrderItemLoader, error) {
	return &supplyOrderItemLoader{
		accountDataSource,
		personDataSource,
		purchaseOrderToSupplyDataSource,
		addressDataSource,
	}, nil
}

func (supOrderItemLoader *supplyOrderItemLoader) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	purchaseOrderToSupply *model.PurchaseOrderToSupplyForSupplyOrderItemInput,
	pickUp *model.DatabaseSupplyOrderItemPickUp,
) (bool, error) {
	purchaseOrderToSupplyLoadedChan := make(chan bool)
	pickUpLoadedChan := make(chan bool)
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
		if pickUp == nil {
			pickUpLoadedChan <- true
			return
		}

		addressLoadedChan := make(chan bool)
		accountLoadedChan := make(chan bool)

		go func() {
			if pickUp.Address == nil {
				addressLoadedChan <- true
				return
			}
			loadedAddress, err := supOrderItemLoader.addressDataSource.GetMongoDataSource().FindByID(
				pickUp.Address.ID,
				session,
			)
			if err != nil {
				errChan <- err
				return
			}
			jsonTemp, _ := json.Marshal(
				map[string]interface{}{
					"Address":  loadedAddress,
					"Status":   model.PickUpStatusAddressNoted,
					"Courier":  model.AccountForSupplyOrderItemInput{},
					"PublicID": "",
				},
			)
			json.Unmarshal(jsonTemp, pickUp)

			addressLoadedChan <- true
		}()

		go func() {
			if pickUp.Courier == nil {
				accountLoadedChan <- true
				return
			}
			account, err := supOrderItemLoader.accountDataSource.GetMongoDataSource().FindByID(
				pickUp.Courier.ID,
				session,
			)
			if err != nil {
				errChan <- err
				return
			}

			jsonTemp, _ := json.Marshal(account)
			json.Unmarshal(jsonTemp, &pickUp.Courier)

			person, err := supOrderItemLoader.personDataSource.GetMongoDataSource().FindByID(
				pickUp.Courier.Person.ID,
				session,
			)
			if err != nil {
				errChan <- err
				return
			}

			jsonTemp, _ = json.Marshal(person)
			json.Unmarshal(jsonTemp, &pickUp.Courier.Person)

			accountLoadedChan <- true
		}()

		for i := 0; i < 2; {
			select {
			case _ = <-addressLoadedChan:
				i++
			case _ = <-accountLoadedChan:
				i++
			}
		}
		pickUpLoadedChan <- true
	}()

	for i := 0; i < 2; {
		select {
		case err := <-errChan:
			return false, err
		case _ = <-purchaseOrderToSupplyLoadedChan:
			i++
		case _ = <-pickUpLoadedChan:
			i++
		}
	}

	return true, nil
}
