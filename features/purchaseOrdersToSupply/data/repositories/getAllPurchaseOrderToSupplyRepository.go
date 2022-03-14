package purchaseordertosupplydomainrepositories

import (
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasepurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/interfaces/sources"
	purchaseordertosupplydomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories"
	purchaseordertosupplydomainrepositorytypes "github.com/horeekaa/backend/features/purchaseOrdersToSupply/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllPurchaseOrderToSupplyRepository struct {
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource
	mongoQueryBuilder               mongodbcorequerybuilderinterfaces.MongoQueryBuilder
	pathIdentity                    string
}

func NewGetAllPurchaseOrderToSupplyRepository(
	purchaseOrderToSupplyDataSource databasepurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (purchaseordertosupplydomainrepositoryinterfaces.GetAllPurchaseOrderToSupplyRepository, error) {
	return &getAllPurchaseOrderToSupplyRepository{
		purchaseOrderToSupplyDataSource,
		mongoQueryBuilder,
		"GetAllPurchaseOrderToSupplyRepository",
	}, nil
}

func (getAllPoToSupplyRepo *getAllPurchaseOrderToSupplyRepository) Execute(
	input purchaseordertosupplydomainrepositorytypes.GetAllPurchaseOrderToSupplyInput,
) ([]*model.PurchaseOrderToSupply, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllPoToSupplyRepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	purchaseOrdersToSupply, err := getAllPoToSupplyRepo.purchaseOrderToSupplyDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getAllPoToSupplyRepo.pathIdentity,
			err,
		)
	}

	return purchaseOrdersToSupply, nil
}
