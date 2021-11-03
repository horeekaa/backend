package purchaseordertosupplydomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasepurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getPurchaseOrderToSupplyRepository struct {
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource
}

func NewGetPurchaseOrderToSupplyRepository(
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource,
) (purchaseordertosupplydomainrepositoryinterfaces.GetPurchaseOrderToSupplyRepository, error) {
	return &getPurchaseOrderToSupplyRepository{
		purchaseOrderToSupplyDataSource,
	}, nil
}

func (getPOToSupplyRepo *getPurchaseOrderToSupplyRepository) Execute(filterFields *model.PurchaseOrderToSupplyFilterFields) (*model.PurchaseOrderToSupply, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	purchaseOrderToSupply, err := getPOToSupplyRepo.purchaseOrderToSupplyDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getPurchaseOrderToSupply",
			err,
		)
	}

	return purchaseOrderToSupply, nil
}
