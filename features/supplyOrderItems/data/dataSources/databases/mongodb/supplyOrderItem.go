package mongodbsupplyorderitemdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbsupplyorderitemdatasourceinterfaces "github.com/horeekaa/backend/features/supplyOrderItems/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type supplyOrderItemDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewSupplyOrderItemDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbsupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSourceMongo, error) {
	basicOperation.SetCollection("supplyorderitems")
	return &supplyOrderItemDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
}

func (supOrderItemDataSourceMongo *supplyOrderItemDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (supOrderItemDataSourceMongo *supplyOrderItemDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.SupplyOrderItem, error) {
	var output model.SupplyOrderItem
	_, err := supOrderItemDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (supOrderItemDataSourceMongo *supplyOrderItemDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.SupplyOrderItem, error) {
	var output model.SupplyOrderItem
	_, err := supOrderItemDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (supOrderItemDataSourceMongo *supplyOrderItemDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.SupplyOrderItem, error) {
	var supplyOrderItems = []*model.SupplyOrderItem{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var supplyOrderItem model.SupplyOrderItem
		if err := cursor.Decode(&supplyOrderItem); err != nil {
			return err
		}
		supplyOrderItems = append(supplyOrderItems, &supplyOrderItem)
		return nil
	}
	_, err := supOrderItemDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return supplyOrderItems, err
}

func (supOrderItemDataSourceMongo *supplyOrderItemDataSourceMongo) Create(input *model.DatabaseCreateSupplyOrderItem, operationOptions *mongodbcoretypes.OperationOptions) (*model.SupplyOrderItem, error) {
	_, err := supOrderItemDataSourceMongo.setDefaultValuesWhenCreate(
		input,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.SupplyOrderItem
	_, err = supOrderItemDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (supOrderItemDataSourceMongo *supplyOrderItemDataSourceMongo) Update(
	updateCriteria map[string]interface{},
	updateData *model.DatabaseUpdateSupplyOrderItem,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.SupplyOrderItem, error) {
	_, err := supOrderItemDataSourceMongo.setDefaultValuesWhenUpdate(
		updateCriteria,
		updateData,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.SupplyOrderItem
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

func (supOrderItemDataSourceMongo *supplyOrderItemDataSourceMongo) setDefaultValuesWhenUpdate(
	inputCriteria map[string]interface{},
	input *model.DatabaseUpdateSupplyOrderItem,
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
			"/supplyOrderItemDataSource/update",
			nil,
		)
	}

	if input.ProposedChanges != nil {
		input.ProposedChanges.UpdatedAt = &currentTime
	}

	return true, nil
}

func (supOrderItemDataSourceMongo *supplyOrderItemDataSourceMongo) setDefaultValuesWhenCreate(
	input *model.DatabaseCreateSupplyOrderItem,
) (bool, error) {
	currentTime := time.Now()
	defaultProposalStatus := model.EntityProposalStatusProposed

	if input.ProposalStatus == nil {
		input.ProposalStatus = &defaultProposalStatus
	}

	input.CreatedAt = &currentTime
	input.UpdatedAt = &currentTime
	if input.ProposedChanges != nil {
		input.ProposedChanges.UpdatedAt = &currentTime
	}

	return true, nil
}
