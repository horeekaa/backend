package mongodbmouitemdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbmouitemdatasourceinterfaces "github.com/horeekaa/backend/features/mouItems/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mouItemDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewMouItemDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbmouitemdatasourceinterfaces.MouItemDataSourceMongo, error) {
	basicOperation.SetCollection("mouitems")
	return &mouItemDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
}

func (mouItemDataSourceMongo *mouItemDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (mouItemDataSourceMongo *mouItemDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.MouItem, error) {
	var output model.MouItem
	_, err := mouItemDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (mouItemDataSourceMongo *mouItemDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.MouItem, error) {
	var output model.MouItem
	_, err := mouItemDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (mouItemDataSourceMongo *mouItemDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.MouItem, error) {
	var mouItems = []*model.MouItem{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var mouItem model.MouItem
		if err := cursor.Decode(&mouItem); err != nil {
			return err
		}
		mouItems = append(mouItems, &mouItem)
		return nil
	}
	_, err := mouItemDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return mouItems, err
}

func (mouItemDataSourceMongo *mouItemDataSourceMongo) Create(input *model.DatabaseCreateMouItem, operationOptions *mongodbcoretypes.OperationOptions) (*model.MouItem, error) {
	_, err := mouItemDataSourceMongo.setDefaultValuesWhenCreate(
		input,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.MouItem
	_, err = mouItemDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (mouItemDataSourceMongo *mouItemDataSourceMongo) Update(
	updateCriteria map[string]interface{},
	updateData *model.DatabaseUpdateMouItem,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.MouItem, error) {
	_, err := mouItemDataSourceMongo.setDefaultValuesWhenUpdate(
		updateCriteria,
		updateData,
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.MouItem
	_, err = mouItemDataSourceMongo.basicOperation.Update(
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

func (mouItemDataSourceMongo *mouItemDataSourceMongo) setDefaultValuesWhenUpdate(
	inputCriteria map[string]interface{},
	input *model.DatabaseUpdateMouItem,
	operationOptions *mongodbcoretypes.OperationOptions,
) (bool, error) {
	currentTime := time.Now()
	existingObject, err := mouItemDataSourceMongo.FindOne(inputCriteria, operationOptions)
	if err != nil {
		return false, err
	}
	if existingObject == nil {
		return false, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.QueryObjectFailed,
			"/mouItemDataSource/update",
			nil,
		)
	}
	input.UpdatedAt = &currentTime

	return true, nil
}

func (mouItemDataSourceMongo *mouItemDataSourceMongo) setDefaultValuesWhenCreate(
	input *model.DatabaseCreateMouItem,
) (bool, error) {
	currentTime := time.Now()
	defaultIsActive := true

	input.CreatedAt = &currentTime
	input.UpdatedAt = &currentTime
	if input.IsActive == nil {
		input.IsActive = &defaultIsActive
	}

	return true, nil
}
