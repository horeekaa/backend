package purchaseorderitemdomainrepositories

import (
	"encoding/json"

	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexceptiontofailure "github.com/horeekaa/backend/core/errors/failures/exceptionToFailure"
	databasepurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/interfaces/sources"
	purchaseorderitemdomainrepositoryinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories"
	purchaseorderitemdomainrepositoryutilityinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/domain/repositories/utils"
	"github.com/horeekaa/backend/model"
)

type createPurchaseOrderItemTransactionComponent struct {
	purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource
	purchaseOrderItemLoader     purchaseorderitemdomainrepositoryutilityinterfaces.PurchaseOrderItemLoader
}

func NewCreatePurchaseOrderItemTransactionComponent(
	purchaseOrderItemDataSource databasepurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSource,
	purchaseOrderItemLoader purchaseorderitemdomainrepositoryutilityinterfaces.PurchaseOrderItemLoader,
) (purchaseorderitemdomainrepositoryinterfaces.CreatePurchaseOrderItemTransactionComponent, error) {
	return &createPurchaseOrderItemTransactionComponent{
		purchaseOrderItemDataSource,
		purchaseOrderItemLoader,
	}, nil
}

func (createPurchaseOrderItemTrx *createPurchaseOrderItemTransactionComponent) PreTransaction(
	createPurchaseOrderItemInput *model.InternalCreatePurchaseOrderItem,
) (*model.InternalCreatePurchaseOrderItem, error) {
	return createPurchaseOrderItemInput, nil
}

func (createPurchaseOrderItemTrx *createPurchaseOrderItemTransactionComponent) TransactionBody(
	session *mongodbcoretypes.OperationOptions,
	createPurchaseOrderItemInput *model.InternalCreatePurchaseOrderItem,
) (*model.PurchaseOrderItem, error) {
	purchaseOrderItemToCreate := &model.DatabaseCreatePurchaseOrderItem{}
	jsonTemp, _ := json.Marshal(createPurchaseOrderItemInput)
	json.Unmarshal(jsonTemp, purchaseOrderItemToCreate)

	_, err := createPurchaseOrderItemTrx.purchaseOrderItemLoader.TransactionBody(
		session,
		purchaseOrderItemToCreate.MouItem,
		purchaseOrderItemToCreate.ProductVariant,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createPurchaseOrderItem",
			err,
		)
	}

	createdPurchaseOrderItem, err := createPurchaseOrderItemTrx.purchaseOrderItemDataSource.GetMongoDataSource().Create(
		purchaseOrderItemToCreate,
		session,
	)
	if err != nil {
		return nil, horeekaacoreexceptiontofailure.ConvertException(
			"/createPurchaseOrderItem",
			err,
		)
	}

	return createdPurchaseOrderItem, nil
}