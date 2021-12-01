package purchaseorderdomainrepositoryutilities

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type purchaseOrderLoader struct {
	mouDataSource          databasemoudatasourceinterfaces.MouDataSource
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource
}

func NewPurchaseOrderLoader(
	mouDataSource databasemoudatasourceinterfaces.MouDataSource,
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
) (purchaseorderdomainrepositoryutilityinterfaces.PurchaseOrderLoader, error) {
	return &purchaseOrderLoader{
		mouDataSource,
		organizationDataSource,
	}, nil
}

func (purcOrderLoader *purchaseOrderLoader) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	mou *model.MouForPurchaseOrderInput,
	organization *model.OrganizationForPurchaseOrderInput,
) (bool, error) {
	mouLoadedChan := make(chan bool)
	organizationLoadedChan := make(chan bool)
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

	for i := 0; i < 2; {
		select {
		case err := <-errChan:
			return false, err
		case _ = <-mouLoadedChan:
			i++
		case _ = <-organizationLoadedChan:
			i++
		}
	}

	return true, nil
}
