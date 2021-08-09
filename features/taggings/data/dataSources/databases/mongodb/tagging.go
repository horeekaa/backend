package mongodbtaggingdatasources

import (
	"time"

	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	mongodbtaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taggingDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
}

func NewTaggingDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbtaggingdatasourceinterfaces.TaggingDataSourceMongo, error) {
	basicOperation.SetCollection("taggings")
	return &taggingDataSourceMongo{
		basicOperation: basicOperation,
	}, nil
}

func (taggingDataSourceMongo *taggingDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (taggingDataSourceMongo *taggingDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tagging, error) {
	var output model.Tagging
	_, err := taggingDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (taggingDataSourceMongo *taggingDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tagging, error) {
	var output model.Tagging
	_, err := taggingDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (taggingDataSourceMongo *taggingDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.Tagging, error) {
	var taggings = []*model.Tagging{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var tagging model.Tagging
		if err := cursor.Decode(&tagging); err != nil {
			return err
		}
		taggings = append(taggings, &tagging)
		return nil
	}
	_, err := taggingDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return taggings, err
}

func (taggingDataSourceMongo *taggingDataSourceMongo) Create(input *model.DatabaseCreateTagging, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tagging, error) {
	defaultedInput, err := taggingDataSourceMongo.setDefaultValues(*input,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesCreateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var outputModel model.Tagging
	_, err = taggingDataSourceMongo.basicOperation.Create(*defaultedInput.CreateTagging, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (taggingDataSourceMongo *taggingDataSourceMongo) Update(ID primitive.ObjectID, updateData *model.DatabaseUpdateTagging, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tagging, error) {
	updateData.ID = ID
	defaultedInput, err := taggingDataSourceMongo.setDefaultValues(*updateData,
		&mongodbcoretypes.DefaultValuesOptions{DefaultValuesType: mongodbcoretypes.DefaultValuesUpdateType},
		operationOptions,
	)
	if err != nil {
		return nil, err
	}

	var output model.Tagging
	_, err = taggingDataSourceMongo.basicOperation.Update(ID, *defaultedInput.UpdateTagging, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

type setTaggingDefaultValuesOutput struct {
	CreateTagging *model.DatabaseCreateTagging
	UpdateTagging *model.DatabaseUpdateTagging
}

func (taggingDataSourceMongo *taggingDataSourceMongo) setDefaultValues(input interface{}, options *mongodbcoretypes.DefaultValuesOptions, operationOptions *mongodbcoretypes.OperationOptions) (*setTaggingDefaultValuesOutput, error) {
	currentTime := time.Now()
	defaultIsActive := true

	if (*options).DefaultValuesType == mongodbcoretypes.DefaultValuesUpdateType {
		updateInput := input.(model.DatabaseUpdateTagging)
		_, err := taggingDataSourceMongo.FindByID(updateInput.ID, operationOptions)
		if err != nil {
			return nil, err
		}
		updateInput.UpdatedAt = &currentTime

		return &setTaggingDefaultValuesOutput{
			UpdateTagging: &updateInput,
		}, nil
	}
	createInput := (input).(model.DatabaseCreateTagging)
	if createInput.IsActive == nil {
		createInput.IsActive = &defaultIsActive
	}
	createInput.CreatedAt = &currentTime
	createInput.UpdatedAt = &currentTime

	return &setTaggingDefaultValuesOutput{
		CreateTagging: &createInput,
	}, nil
}
