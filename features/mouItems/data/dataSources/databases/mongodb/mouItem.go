package mongodbmouitemdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
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
	defaultedInput, err := mouItemDataSourceMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.MouItem
	_, err = mouItemDataSourceMongo.basicOperation.Create(*defaultedInput.CreateMouItem, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (mouItemDataSourceMongo *mouItemDataSourceMongo) Update(ID primitive.ObjectID, updateData *model.DatabaseUpdateMouItem, operationOptions *mongodbcoretypes.OperationOptions) (*model.MouItem, error) {
	updateData.ID = ID
	defaultedInput, err := mouItemDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.MouItem
	_, err = mouItemDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdateMouItem, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

type setMouItemDefaultValuesOutput struct {
	CreateMouItem *model.DatabaseCreateMouItem
	UpdateMouItem *model.DatabaseUpdateMouItem
}

func (mouItemDataSourceMongo *mouItemDataSourceMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setMouItemDefaultValuesOutput, error) {
	currentTime := time.Now()
	defaultIsActive := true

	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		updateInput := input.(model.DatabaseUpdateMouItem)
		_, err := mouItemDataSourceMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}
		updateInput.UpdatedAt = &currentTime

		return &setMouItemDefaultValuesOutput{
			UpdateMouItem: &updateInput,
		}, nil
	}
	createInput := (input).(model.DatabaseCreateMouItem)
	createInput.CreatedAt = &currentTime
	createInput.UpdatedAt = &currentTime
	if createInput.IsActive == nil {
		createInput.IsActive = &defaultIsActive
	}

	return &setMouItemDefaultValuesOutput{
		CreateMouItem: &createInput,
	}, nil
}
