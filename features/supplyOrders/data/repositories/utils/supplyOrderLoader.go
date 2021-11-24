package supplyorderdomainrepositoryutilities

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databaseaddressdatasourceinterfaces "github.com/horeekaa/backend/features/addresses/data/dataSources/databases/interfaces/sources"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type supplyOrderLoader struct {
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource
	addressDataSource      databaseaddressdatasourceinterfaces.AddressDataSource
}

func NewSupplyOrderLoader(
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
	addressDataSource databaseaddressdatasourceinterfaces.AddressDataSource,
) (supplyorderdomainrepositoryutilityinterfaces.SupplyOrderLoader, error) {
	return &supplyOrderLoader{
		organizationDataSource,
		addressDataSource,
	}, nil
}

func (purcOrderLoader *supplyOrderLoader) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	organization *model.OrganizationForSupplyOrderInput,
	address *model.AddressForSupplyOrderInput,
) (bool, error) {
	organizationLoadedChan := make(chan bool)
	addressLoadedChan := make(chan bool)
	errChan := make(chan error)

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

	for i := 0; i < 2; {
		select {
		case err := <-errChan:
			return false, err
		case _ = <-organizationLoadedChan:
			i++
		case _ = <-addressLoadedChan:
			i++
		}
	}

	return true, nil
}
