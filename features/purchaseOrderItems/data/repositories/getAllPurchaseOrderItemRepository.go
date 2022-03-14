package purchaseorderitemdomainrepositories

import (
	mongodbcorequerybuilderinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/queryBuilders"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	purchaseorderitemdomainrepositorytypes "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories/types"
	"github.com/horeekaa/backend/model"
)

type getAllPurchaseOrderItemRepository struct {
	purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource
	mongoQueryBuilder           mongodbcorequerybuilderinterfaces.MongoQueryBuilder
	pathIdentity                string
}

func NewGetAllPurchaseOrderItemRepository(
	purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
	mongoQueryBuilder mongodbcorequerybuilderinterfaces.MongoQueryBuilder,
) (purchaseorderitemdomainrepositoryinterfaces.GetAllPurchaseOrderItemRepository, error) {
	return &getAllPurchaseOrderItemRepository{
		purchaseOrderItemDataSource,
		mongoQueryBuilder,
		"GetAllPurchaseOrderItemRepository",
	}, nil
}

func (getAllPOItemRepo *getAllPurchaseOrderItemRepository) Execute(
	input purchaseorderitemdomainrepositorytypes.GetAllPurchaseOrderItemInput,
) ([]*model.PurchaseOrderItem, error) {
	filterFieldsMap := map[string]interface{}{}
	getAllPOItemRepo.mongoQueryBuilder.Execute(
		"",
		input.FilterFields,
		&filterFieldsMap,
	)

	mongoPagination := (mongodbcoretypes.PaginationOptions)(*input.PaginationOpt)

	purchaseOrderItems, err := getAllPOItemRepo.purchaseOrderItemDataSource.GetMongoDataSource().Find(
		filterFieldsMap,
		&mongoPagination,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getAllPOItemRepo.pathIdentity,
			err,
		)
	}

	return purchaseOrderItems, nil
}
