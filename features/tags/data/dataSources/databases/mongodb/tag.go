package mongodbtagdatasources

import (
	mongodbcoreoperationinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/operations"
	mongodbcorewrapperinterfaces "github.com/horeekaa/backend/core/databaseClient/mongodb/interfaces/wrappers"
	mongodbcoretypes "github.com/horeekaa/backend/core/databaseClient/mongodb/types"
	horeekaacoreexception "github.com/horeekaa/backend/core/errors/exceptions"
	horeekaacoreexceptionenums "github.com/horeekaa/backend/core/errors/exceptions/enums"
	mongodbtagdatasourceinterfaces "github.com/horeekaa/backend/features/tags/data/dataSources/databases/mongodb/interfaces"
	"github.com/horeekaa/backend/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type tagDataSourceMongo struct {
	basicOperation mongodbcoreoperationinterfaces.BasicOperation
	pathIdentity   string
}

func NewTagDataSourceMongo(basicOperation mongodbcoreoperationinterfaces.BasicOperation) (mongodbtagdatasourceinterfaces.TagDataSourceMongo, error) {
	basicOperation.SetCollection("tags")
	return &tagDataSourceMongo{
		basicOperation: basicOperation,
		pathIdentity:   "TagDataSource",
	}, nil
}

func (tagDataSourceMongo *tagDataSourceMongo) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (tagDataSourceMongo *tagDataSourceMongo) FindByID(ID primitive.ObjectID, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tag, error) {
	var output model.Tag
	_, err := tagDataSourceMongo.basicOperation.FindByID(ID, &output, operationOptions)
	if err != nil {
		return nil, err
	}

	return &output, nil
}

func (tagDataSourceMongo *tagDataSourceMongo) FindOne(query map[string]interface{}, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tag, error) {
	var output model.Tag
	_, err := tagDataSourceMongo.basicOperation.FindOne(query, &output, operationOptions)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &output, err
}

func (tagDataSourceMongo *tagDataSourceMongo) Find(
	query map[string]interface{},
	paginationOpts *mongodbcoretypes.PaginationOptions,
	operationOptions *mongodbcoretypes.OperationOptions,
) ([]*model.Tag, error) {
	var tags = []*model.Tag{}
	appendingFn := func(cursor mongodbcorewrapperinterfaces.MongoCursor) error {
		var tag model.Tag
		if err := cursor.Decode(&tag); err != nil {
			return err
		}
		tags = append(tags, &tag)
		return nil
	}
	_, err := tagDataSourceMongo.basicOperation.Find(query, paginationOpts, appendingFn, operationOptions)
	if err != nil {
		return nil, err
	}

	return tags, err
}

func (tagDataSourceMongo *tagDataSourceMongo) Create(input *model.DatabaseCreateTag, operationOptions *mongodbcoretypes.OperationOptions) (*model.Tag, error) {
	var outputModel model.Tag
	_, err := tagDataSourceMongo.basicOperation.Create(input, &outputModel, operationOptions)
	if err != nil {
		return nil, err
	}

	return &outputModel, err
}

func (tagDataSourceMongo *tagDataSourceMongo) Update(
	updateCriteria map[string]interface{},
	updateData *model.DatabaseUpdateTag,
	operationOptions *mongodbcoretypes.OperationOptions,
) (*model.Tag, error) {
	existingObject, err := tagDataSourceMongo.FindOne(updateCriteria, operationOptions)
	if err != nil {
		return nil, err
	}
	if existingObject == nil {
		return nil, horeekaacoreexception.NewExceptionObject(
			horeekaacoreexceptionenums.NoUpdatableObjectFound,
			tagDataSourceMongo.pathIdentity,
			nil,
		)
	}

	var output model.Tag
	_, err = tagDataSourceMongo.basicOperation.Update(
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
