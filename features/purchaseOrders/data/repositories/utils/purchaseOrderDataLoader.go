package purchaseorderdomainrepositoryutilities

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type purchaseOrderLoader struct {
	mouDataSource          databasemoudatasourceinterfaces.MouDataSource
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource
	addressDataSource      databaseaddressdatasourceinterfaces.AddressDataSource
}

func NewPurchaseOrderLoader(
	mouDataSource databasemoudatasourceinterfaces.MouDataSource,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
) (purchaseorderdomainrepositoryutilityinterfaces.PurchaseOrderLoader, error) {
	return &purchaseOrderLoader{
		mouDataSource,
		organizationDataSource,
		addressDataSource,
	}, nil
}

func (purcOrderLoader *purchaseOrderLoader) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	mou *model.MouForPurchaseOrderInput,
	organization *model.OrganizationForPurchaseOrderInput,
	address *model.AddressForPurchaseOrderInput,
) (bool, error) {
	mouLoadedChan := make(chan bool)
	organizationLoadedChan := make(chan bool)
	addressLoadedChan := make(chan bool)
	errChan := make(chan error)
	go func() {
		if mou == nil {
			mouLoadedChan <- true
			return
		}

		loadedMou, err := purcOrderLoader.mouDataSource.GetMongoDataSource().FindByID(
			mou.ID,
			session,
		)
		if err != nil {
			errChan <- err
			return
		}
		jsonTemp, _ := json.Marshal(loadedMou)
		json.Unmarshal(jsonTemp, mou)

		mouLoadedChan <- true
	}()

	go func() {
		if organization == nil {
			organizationLoadedChan <- true
			return
		}

		loadedOrganization, err := purcOrderLoader.organizationDataSource.GetMongoDataSource().FindByID(
			organization.ID,
			session,
		)
		if err != nil {
			errChan <- err
			return
		}
		jsonTemp, _ := json.Marshal(loadedOrganization)
		json.Unmarshal(jsonTemp, organization)

		organizationLoadedChan <- true
	}()

	go func() {
		if address == nil {
			addressLoadedChan <- true
			return
		}

		loadedAddress, err := purcOrderLoader.addressDataSource.GetMongoDataSource().FindByID(
			address.ID,
			session,
		)
		if err != nil {
			errChan <- err
			return
		}
		jsonTemp, _ := json.Marshal(loadedAddress)
		json.Unmarshal(jsonTemp, address)

		addressLoadedChan <- true
	}()

	for i := 0; i < 3; {
		select {
		case err := <-errChan:
			return false, err
		case _ = <-mouLoadedChan:
			i++
		case _ = <-organizationLoadedChan:
			i++
		case _ = <-addressLoadedChan:
			i++
		}
	}

	return true, nil
}
