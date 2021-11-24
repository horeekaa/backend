package mongodbsupplyorderdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbsupplyorderdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrders/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type supplyOrderDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewSupplyOrderDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbsupplyorderdatasourceinterfaces.SupplyOrderDataSourceMongo, error) {
	basicOperation.SetCollection("supplyorders")
	return &supplyOrderDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
}

func (supOrderItemDataSourceMongo *supplyOrderDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (supOrderItemDataSourceMongo *supplyOrderDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.SupplyOrder, error) {
	var output model.SupplyOrder
	_, err := supOrderItemDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (supOrderItemDataSourceMongo *supplyOrderDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.SupplyOrder, error) {
	var output model.SupplyOrder
	_, err := supOrderItemDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (supOrderItemDataSourceMongo *supplyOrderDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.SupplyOrder, error) {
	var supplyOrders = []*model.SupplyOrder{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var supplyOrder model.SupplyOrder
		if err := cursor.Decode(&supplyOrder); err != nil {
			return err
		}
		supplyOrders = append(supplyOrders, &supplyOrder)
		return nil
	}
	_, err := supOrderItemDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return supplyOrders, err
}

func (supOrderItemDataSourceMongo *supplyOrderDataSourceMongo) Create(input *model.DatabaseCreateSupplyOrder, operationOptions *mongodbcoretypes.OperationOptions) (*model.SupplyOrder, error) {
	_, err := supOrderItemDataSourceMongo.setDefaultValuesWhenCreate(
		input,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.SupplyOrder
	_, err = supOrderItemDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (supOrderItemDataSourceMongo *supplyOrderDataSourceMongo) Update(
	updateCriteria map[string]interface{},
	updateData *model.DatabaseUpdateSupplyOrder,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.SupplyOrder, error) {
	_, err := supOrderItemDataSourceMongo.setDefaultValuesWhenUpdate(
		updateCriteria,
		updateData,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.SupplyOrder
	_, err = supOrderItemDataSourceMongo.basicOperation.Update(
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

func (supOrderItemDataSourceMongo *supplyOrderDataSourceMongo) setDefaultValuesWhenUpdate(
	inputCriteria map[string]interface{},
	input *model.DatabaseUpdateSupplyOrder,
	operationOptions *mongodbcoretypes.OperationOptions,
) (bool, error) {
	currentTime := time.Now()
	existingObject, err := supOrderItemDataSourceMongo.FindOne(inputCriteria, operationOptions)
	if err != nil {
		return false, err
	}
	if existingObject == nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.QueryObjectFailed,
			"/supplyOrderDataSource/update",
			nil,
		)
	}

	if input.ProposedChanges != nil {
		input.ProposedChanges.UpdatedAt = &currentTime
	}

	return true, nil
}

func (supOrderItemDataSourceMongo *supplyOrderDataSourceMongo) setDefaultValuesWhenCreate(
	input *model.DatabaseCreateSupplyOrder,
) (bool, error) {
	currentTime := time.Now()
	defaultProposalStatus := model.EntityProposalStatusProposed

	if input.ProposalStatus == nil {
		input.ProposalStatus = &defaultProposalStatus
	}

	input.CreatedAt = currentTime
	input.UpdatedAt = currentTime
	if input.ProposedChanges != nil {
		input.ProposedChanges.UpdatedAt = &currentTime
	}

	return true, nil
}
