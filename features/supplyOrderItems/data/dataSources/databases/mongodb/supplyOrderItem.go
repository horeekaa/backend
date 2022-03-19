package mongodbsupplyorderitemdatasources

import (
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
	pathIdentity   string
}

func NewSupplyOrderItemDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbsupplyorderitemdatasourceinterfaces.SupplyOrderItemDataSourceMongo, error) {
	basicOperation.SetCollection("supplyorderitems")
	return &supplyOrderItemDataSourceMongo{
		basicOperation: basicOperation,
		pathIdentity:   "SupplyOrderItemDataSource",
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
	var outputModel model.SupplyOrderItem
	_, err := supOrderItemDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
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
	existingObject, err := supOrderItemDataSourceMongo.FindOne(updateCriteria, operationOptions)
	if err != nil {
		return nil, err
	}
	if existingObject == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.NoUpdatableObjectFound,
			supOrderItemDataSourceMongo.pathIdentity,
			nil,
		)
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
