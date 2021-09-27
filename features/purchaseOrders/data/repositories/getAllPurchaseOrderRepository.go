package purchaseorderdomainrepositories

import (
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasepurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/interfaces/sources"
	purchaseorderdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories"
	purchaseorderdomainrepositorytypes "github.com/horeekaa/backend/features/purchaseOrders/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllPurchaseOrderRepository struct {
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource
	mongoQueryBuilder       mongodbcorequerybuilderinterfaces.MongoQueryBuilder
}

func NewGetAllPurchaseOrderRepository(
	purchaseOrderDataSource databasepurchaseorderdatasourceinterfaces.PurchaseOrderDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (purchaseorderdomainrepositoryinterfaces.GetAllPurchaseOrderRepository, error) {
	return &getAllPurchaseOrderRepository{
		purchaseOrderDataSource,
		mongoQueryBuilder,
	}, nil
}

func (getAllPORepo *getAllPurchaseOrderRepository) Execute(
	input purchaseorderdomainrepositorytypes.GetAllPurchaseOrderInput,
) ([]*model.PurchaseOrder, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllPORepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	purchaseOrders, err := getAllPORepo.purchaseOrderDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/getAllPurchaseOrder",
			err,
		)
	}

	return purchaseOrders, nil
}
