package supplyorderdomainrepositories

import (
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasesupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/interfaces/sources"
	supplyorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrders/domain/repositories"
	supplyorderdomainrepositorytypes "github.com/horeekaa/backend/features/supplyOrders/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllSupplyOrderRepository struct {
	supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource
	mongoQueryBuilder     mongodbcorequerybuilderinterfaces.MongoQueryBuilder
	pathIdentity          string
}

func NewGetAllSupplyOrderRepository(
	supplyOrderDataSource databasesupplyorderdatasourceinterfaces.SupplyOrderDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (supplyorderdomainrepositoryinterfaces.GetAllSupplyOrderRepository, error) {
	return &getAllSupplyOrderRepository{
		supplyOrderDataSource,
		mongoQueryBuilder,
		"GetAllSupplyOrderRepository",
	}, nil
}

func (getAllSORepo *getAllSupplyOrderRepository) Execute(
	input supplyorderdomainrepositorytypes.GetAllSupplyOrderInput,
) ([]*model.SupplyOrder, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllSORepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	supplyOrders, err := getAllSORepo.supplyOrderDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getAllSORepo.pathIdentity,
			err,
		)
	}

	return supplyOrders, nil
}
