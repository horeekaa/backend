package mongodbpurchaseordertosupplydatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbpurchaseordertosupplydatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrdersToSupply/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type purchaseOrderToSupplyDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
	pathIdentity   string
}

func NewPurchaseOrderToSupplyDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbpurchaseordertosupplydatasourceinterfaces.PurchaseOrderToSupplyDataSourceMongo, error) {
	basicOperation.SetCollection("purchaseorderstosupply")
	return &purchaseOrderToSupplyDataSourceMongo{
		basicOperation: basicOperation,
		pathIdentity:   "PurchaseOrderToSupplyDataSource",
	}, nil
}

func (purcOrderSupplyDataSourceMongo *purchaseOrderToSupplyDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrderToSupply, error) {
	var output model.PurchaseOrderToSupply
	_, err := purcOrderSupplyDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (purcOrderSupplyDataSourceMongo *purchaseOrderToSupplyDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrderToSupply, error) {
	var output model.PurchaseOrderToSupply
	_, err := purcOrderSupplyDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (purcOrderSupplyDataSourceMongo *purchaseOrderToSupplyDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.PurchaseOrderToSupply, error) {
	var purchaseOrdersToSupply = []*model.PurchaseOrderToSupply{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var purchaseOrderToSupply model.PurchaseOrderToSupply
		if err := cursor.Decode(&purchaseOrderToSupply); err != nil {
			return err
		}
		purchaseOrdersToSupply = append(purchaseOrdersToSupply, &purchaseOrderToSupply)
		return nil
	}
	_, err := purcOrderSupplyDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return purchaseOrdersToSupply, err
}

func (purcOrderSupplyDataSourceMongo *purchaseOrderToSupplyDataSourceMongo) Create(input *model.DatabaseCreatePurchaseOrderToSupply, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrderToSupply, error) {
	_, err := purcOrderSupplyDataSourceMongo.setDefaultValuesWhenCreate(
		input,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.PurchaseOrderToSupply
	_, err = purcOrderSupplyDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (purcOrderSupplyDataSourceMongo *purchaseOrderToSupplyDataSourceMongo) Update(
	updateCriteria map[string]interface{},
	updateData *model.DatabaseUpdatePurchaseOrderToSupply,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.PurchaseOrderToSupply, error) {
	_, err := purcOrderSupplyDataSourceMongo.setDefaultValuesWhenUpdate(
		updateCriteria,
		updateData,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.PurchaseOrderToSupply
	_, err = purcOrderSupplyDataSourceMongo.basicOperation.Update(
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

func (purcOrderSupplyDataSourceMongo *purchaseOrderToSupplyDataSourceMongo) setDefaultValuesWhenUpdate(
	inputCriteria map[string]interface{},
	input *model.DatabaseUpdatePurchaseOrderToSupply,
	operationOptions *mongodbcoretypes.OperationOptions,
) (bool, error) {
	existingObject, err := purcOrderSupplyDataSourceMongo.FindOne(inputCriteria, operationOptions)
	if err != nil {
		return false, err
	}
	if existingObject == nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.NoUpdatableObjectFound,
			purcOrderSupplyDataSourceMongo.pathIdentity,
			nil,
		)
	}

	return true, nil
}

func (purcOrderSupplyDataSourceMongo *purchaseOrderToSupplyDataSourceMongo) setDefaultValuesWhenCreate(
	input *model.DatabaseCreatePurchaseOrderToSupply,
) (bool, error) {
	currentTime := time.Now().UTC()
	defaultStatus := model.PurchaseOrderToSupplyStatusCummulating

	input.CreatedAt = currentTime
	input.UpdatedAt = currentTime
	if input.Status == nil {
		input.Status = &defaultStatus
	}

	return true, nil
}
