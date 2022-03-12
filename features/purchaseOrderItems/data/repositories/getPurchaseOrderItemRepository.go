package purchaseorderitemdomainrepositories

import (
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson"
)

type getPurchaseOrderItemRepository struct {
	purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource
	pathIdentity                string
}

func NewGetPurchaseOrderItemRepository(
	purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
) (purchaseorderitemdomainrepositoryinterfaces.GetPurchaseOrderItemRepository, error) {
	return &getPurchaseOrderItemRepository{
		purchaseOrderItemDataSource,
		"GetPurchaseOrderItemRepository",
	}, nil
}

func (getPurchaseOrderItemRepo *getPurchaseOrderItemRepository) Execute(filterFields *model.PurchaseOrderItemFilterFields) (*model.PurchaseOrderItem, error) {
	if filterFields == nil {
		return nil, nil
	}

	var filterFieldsMap map[string]interface{}
	data, _ := bson.Marshal(filterFields)
	bson.Unmarshal(data, &filterFieldsMap)

	purchaseOrderItem, err := getPurchaseOrderItemRepo.purchaseOrderItemDataSource.GetMongoDataSource().FindOne(
		filterFieldsMap,
		&mongodbcoretypes.OperationOptions{},
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			getPurchaseOrderItemRepo.pathIdentity,
			err,
		)
	}

	return purchaseOrderItem, nil
}
