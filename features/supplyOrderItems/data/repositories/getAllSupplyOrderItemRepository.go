package supplyorderitemdomainrepositories

import (
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasesupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/interfaces/sources"
	supplyorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories"
	supplyorderitemdomainrepositorytypes "github.com/horeekaa/backend/features/supplyOrderItems/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllSupplyOrderItemRepository struct {
	supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource
	mongoQueryBuilder         mongodbcorequerybuilderinterfaces.MongoQueryBuilder
}

func NewGetAllSupplyOrderItemRepository(
	supplyOrderItemDataSource databasesupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (supplyorderitemdomainrepositoryinterfaces.GetAllSupplyOrderItemRepository, error) {
	return &getAllSupplyOrderItemRepository{
		supplyOrderItemDataSource,
		mongoQueryBuilder,
	}, nil
}

func (getAllSOItemRepo *getAllSupplyOrderItemRepository) Execute(
	input supplyorderitemdomainrepositorytypes.GetAllSupplyOrderItemInput,
) ([]*model.SupplyOrderItem, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllSOItemRepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	supplyOrderItems, err := getAllSOItemRepo.supplyOrderItemDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getAllSupplyOrderItem",
			err,
		)
	}

	return supplyOrderItems, nil
}
