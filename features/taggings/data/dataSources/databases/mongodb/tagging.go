package mongodbtaggingdatasources

import (
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbtaggingdatasourceinterfaces "github.com/horeekaa/backend/features/taggings/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taggingDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
	pathIdentity   string
}

func NewTaggingDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbtaggingdatasourceinterfaces.TaggingDataSourceMongo, error) {
	basicOperation.SetCollection("taggings")
	return &taggingDataSourceMongo{
		basicOperation: basicOperation,
		pathIdentity:   "TaggingDataSource",
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

func (taggingDataSourceMongo *taggingDataSourceMongo) Create(
	input *model.DatabaseCreateTagging,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.Tagging, error) {
	var outputModel model.Tagging
	_, err := taggingDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (taggingDataSourceMongo *taggingDataSourceMongo) Update(
	updateCriteria map[string]interface{},
	updateData *model.DatabaseUpdateTagging,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.Tagging, error) {
	existingObject, err := taggingDataSourceMongo.FindOne(updateCriteria, operationOptions)
	if err != nil {
		return nil, err
	}
	if existingObject == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.NoUpdatableObjectFound,
			taggingDataSourceMongo.pathIdentity,
			nil,
		)
	}

	var output model.Tagging
	_, err = taggingDataSourceMongo.basicOperation.Update(
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
