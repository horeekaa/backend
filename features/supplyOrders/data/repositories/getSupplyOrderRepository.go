package supplyorderdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasesupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getSupplyOrderRepository struct {
	supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource
	pathIdentity          string
}

func NewGetSupplyOrderRepository(
	supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource,
) (supplyorderdomainrepositoryinterfaces.GetSupplyOrderRepository, error) {
	return &getSupplyOrderRepository{
		supplyOrderDataSource,
		"GetSupplyOrderRepository",
	}, nil
}

func (getSupplyOrderRepo *getSupplyOrderRepository) Execute(filterFields *model.SupplyOrderFilterFields) (*model.SupplyOrder, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	supplyOrder, err := getSupplyOrderRepo.supplyOrderDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getSupplyOrderRepo.pathIdentity,
			err,
		)
	}

	return supplyOrder, nil
}
