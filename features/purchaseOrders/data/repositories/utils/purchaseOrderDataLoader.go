package purchaseorderdomainrepositoryutilities

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	databasemoudatasourceinterfaces "github.com/horeekaa/backend/features/mous/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type purchaseOrderLoader struct {
	mouDataSource databasemoudatasourceinterfaces.MouDataSource
}

func NewPurchaseOrderLoader(
	mouDataSource databasemoudatasourceinterfaces.MouDataSource,
) (purchaseorderdomainrepositoryutilityinterfaces.PurchaseOrderLoader, error) {
	return &purchaseOrderLoader{
		mouDataSource,
	}, nil
}

func (purcOrderLoader *purchaseOrderLoader) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	mou *model.MouForPurchaseOrderInput,
) (bool, error) {
	if mou == nil {
		return true, nil
	}

	loadedMou, err := purcOrderLoader.mouDataSource.GetMongoDataSource().FindByID(
		mou.ID,
		session,
	)
	if err != nil {
		return false, err
	}

	jsonTemp, _ := json.Marshal(loadedMou)
	json.Unmarshal(jsonTemp, mou)

	return true, nil
}
