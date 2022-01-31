package mongodbpurchaseorderdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbpurchaseorderdatasourceinterfaces "github.com/horeekaa/backend/features/purchaseOrders/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type purchaseOrderDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewPurchaseOrderDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbpurchaseorderdatasourceinterfaces.PurchaseOrderDataSourceMongo, error) {
	basicOperation.SetCollection("purchaseorders")
	return &purchaseOrderDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
}

func (purcOrderDataSourceMongo *purchaseOrderDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (purcOrderDataSourceMongo *purchaseOrderDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrder, error) {
	var output model.PurchaseOrder
	_, err := purcOrderDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (purcOrderDataSourceMongo *purchaseOrderDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrder, error) {
	var output model.PurchaseOrder
	_, err := purcOrderDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (purcOrderDataSourceMongo *purchaseOrderDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.PurchaseOrder, error) {
	var purchaseOrders = []*model.PurchaseOrder{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var purchaseOrder model.PurchaseOrder
		if err := cursor.Decode(&purchaseOrder); err != nil {
			return err
		}
		purchaseOrders = append(purchaseOrders, &purchaseOrder)
		return nil
	}
	_, err := purcOrderDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return purchaseOrders, err
}

func (purcOrderDataSourceMongo *purchaseOrderDataSourceMongo) Create(input *model.DatabaseCreatePurchaseOrder, operationOptions *mongodbcoretypes.OperationOptions) (*model.PurchaseOrder, error) {
	_, err := purcOrderDataSourceMongo.setDefaultValuesWhenCreate(
		input,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.PurchaseOrder
	_, err = purcOrderDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (purcOrderDataSourceMongo *purchaseOrderDataSourceMongo) Update(
	updateCriteria map[string]interface{},
	updateData *model.DatabaseUpdatePurchaseOrder,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.PurchaseOrder, error) {
	_, err := purcOrderDataSourceMongo.setDefaultValuesWhenUpdate(
		updateCriteria,
		updateData,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.PurchaseOrder
	_, err = purcOrderDataSourceMongo.basicOperation.Update(
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

func (purchaseOrderDataSourceMongo *purchaseOrderDataSourceMongo) UpdateAll(
	updateCriteria map[string]interface{},
	updateData *model.DatabaseUpdatePurchaseOrder,
	operationOptions *mongodbcoretypes.OperationOptions,
) (bool, error) {
	currentTime := time.Now()
	updateData.UpdatedAt = &currentTime

	_, err := purchaseOrderDataSourceMongo.basicOperation.UpdateAll(
		updateCriteria,
		map[string]interface{}{
			"$set": updateData,
		},
		operationOptions,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (purcOrderDataSourceMongo *purchaseOrderDataSourceMongo) setDefaultValuesWhenUpdate(
	inputCriteria map[string]interface{},
	input *model.DatabaseUpdatePurchaseOrder,
	operationOptions *mongodbcoretypes.OperationOptions,
) (bool, error) {
	currentTime := time.Now()
	existingObject, err := purcOrderDataSourceMongo.FindOne(inputCriteria, operationOptions)
	if err != nil {
		return false, err
	}
	if existingObject == nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.QueryObjectFailed,
			"/purchaseOrderDataSource/update",
			nil,
		)
	}

	if input.ProposedChanges != nil {
		input.ProposedChanges.UpdatedAt = &currentTime
	}

	return true, nil
}

func (purcOrderDataSourceMongo *purchaseOrderDataSourceMongo) setDefaultValuesWhenCreate(
	input *model.DatabaseCreatePurchaseOrder,
) (bool, error) {
	currentTime := time.Now()
	defaultProposalStatus := model.EntityProposalStatusProposed
	defaultStatus := model.PurchaseOrderStatusOpen

	if input.ProposalStatus == nil {
		input.ProposalStatus = &defaultProposalStatus
	}
	if input.Status == nil {
		input.Status = &defaultStatus
	}
	input.CreatedAt = &currentTime
	input.UpdatedAt = &currentTime
	if input.ProposedChanges != nil {
		input.ProposedChanges.UpdatedAt = &currentTime
	}

	return true, nil
}
