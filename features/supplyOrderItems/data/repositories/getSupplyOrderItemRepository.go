package supplyorderitemdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getSupplyOrderItemRepository struct {
	supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource
}

func NewGetSupplyOrderItemRepository(
	supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource,
) (supplyorderitemdomainrepositoryinterfaces.GetSupplyOrderItemRepository, error) {
	return &getSupplyOrderItemRepository{
		supplyOrderItemDataSource,
	}, nil
}

func (getSupplyOrderItemRepo *getSupplyOrderItemRepository) Execute(filterFields *model.SupplyOrderItemFilterFields) (*model.SupplyOrderItem, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	supplyOrderItem, err := getSupplyOrderItemRepo.supplyOrderItemDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getSupplyOrderItem",
			err,
		)
	}

	return supplyOrderItem, nil
}
