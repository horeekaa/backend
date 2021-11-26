package supplyorderdomainrepositoryutilities

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databaseorganizationdatasourceinterfaces "github.com/horeekaa/backend/features/organizations/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type supplyOrderLoader struct {
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource
}

func NewSupplyOrderLoader(
	organizationDataSource databaseorganizationdatasourceinterfaces.OrganizationDataSource,
) (supplyorderdomainrepositoryutilityinterfaces.SupplyOrderLoader, error) {
	return &supplyOrderLoader{
		organizationDataSource,
	}, nil
}

func (supOrderLoader *supplyOrderLoader) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	organization *model.OrganizationForSupplyOrderInput,
) (bool, error) {
	if organization == nil {
		return true, nil
	}

	loadedOrganization, err := supOrderLoader.organizationDataSource.GetMongoDataSource().FindByID(
		organization.ID,
		session,
	)
	if err != nil {
		return false, err
	}
	jsonTemp, _ := json.Marshal(loadedOrganization)
	json.Unmarshal(jsonTemp, organization)

	return true, nil
}
