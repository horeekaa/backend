package mongodbpurchaseorderItemdatasources

import (
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbpurchaseorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrderItems/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type purchaseOrderItemDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
	pathIdentity   string
}

func NewPurchaseOrderItemDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbpurchaseorderitemdatasourceinterfaces.PurchaseOrderItemDataSourceMongo, error) {
	basicOperation.SetCollection("purchaseorderitems")
	return &purchaseOrderItemDataSourceMongo{
		basicOperation: basicOperation,
		pathIdentity:   "PurchaseOrderItemDataSource",
	}, nil
}

func (purcOrderItemDataSourceMongo *purchaseOrderItemDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (purcOrderItemDataSourceMongo *purchaseOrderItemDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrderItem, error) {
	var output model.PurchaseOrderItem
	_, err := purcOrderItemDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (purcOrderItemDataSourceMongo *purchaseOrderItemDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrderItem, error) {
	var output model.PurchaseOrderItem
	_, err := purcOrderItemDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (purcOrderItemDataSourceMongo *purchaseOrderItemDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.PurchaseOrderItem, error) {
	var purchaseOrderItems = []*model.PurchaseOrderItem{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var purchaseOrderItem model.PurchaseOrderItem
		if err := cursor.Decode(&purchaseOrderItem); err != nil {
			return err
		}
		purchaseOrderItems = append(purchaseOrderItems, &purchaseOrderItem)
		return nil
	}
	_, err := purcOrderItemDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return purchaseOrderItems, err
}

func (purcOrderItemDataSourceMongo *purchaseOrderItemDataSourceMongo) Create(input *model.DatabaseCreatePurchaseOrderItem, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrderItem, error) {
	var outputModel model.PurchaseOrderItem
	_, err := purcOrderItemDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (purcOrderItemDataSourceMongo *purchaseOrderItemDataSourceMongo) Update(
	updateCriteria map[string]interface{},
	updateData *model.DatabaseUpdatePurchaseOrderItem,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.PurchaseOrderItem, error) {
	existingObject, err := purcOrderItemDataSourceMongo.FindOne(updateCriteria, operationOptions)
	if err != nil {
		return nil, err
	}
	if existingObject == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.NoUpdatableObjectFound,
			purcOrderItemDataSourceMongo.pathIdentity,
			nil,
		)
	}

	var output model.PurchaseOrderItem
	_, err = purcOrderItemDataSourceMongo.basicOperation.Update(
		updateCriteria,
		map[string]interface{}{
			"$set": updateData,
		},
		&output,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
