package purchaseorderdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getPurchaseOrderRepository struct {
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
}

func NewGetPurchaseOrderRepository(
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
) (purchaseorderdomainrepositoryinterfaces.GetPurchaseOrderRepository, error) {
	return &getPurchaseOrderRepository{
		purchaseOrderDataSource,
	}, nil
}

func (getPORepo *getPurchaseOrderRepository) Execute(filterFields *model.PurchaseOrderFilterFields) (*model.PurchaseOrder, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	purchaseOrder, err := getPORepo.purchaseOrderDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getPurchaseOrder",
			err,
		)
	}

	return purchaseOrder, nil
}
